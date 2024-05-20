# Gravwell CLI

A reimplementation of the Gravwell CLI, built on the fantastic Bubble Tea and Cobra libraries.

# Design

## Terminology
Bubble Tea has the tea.Model interface that must be implemented by a model struct of our own, Bubbles.TextInput is a tea.Model under the hood. Cobra is composed of cobra.Commands and Bubble Tea drives its I/O via tea.Cmds.

So we are using our own terminology so we can avoid further homonyms. 

The Bubble Tea model struct, our controller, is *Mother*.

Commands that can be invoked (interactively or from a script), such as `search`, are *Actions*.

Menus and child menus, commands that require further input, such as `admin`, are *Navs*.

## Wrapping the Cobra Command Tree

I tried really hard to rely solely on Cobra's underlying tree navigation (via `.Parent()`, `.Root()`, `.Commands()`). It has all the navigational features we need and sticking to it allows Cobra to implicitly generate the entire non-interactive interface to the CLI (perfect for script compatability).

The problem comes in when we need to supplant Mother's `.Update()` and `.View()` with that of the active Action. We know a cobra.Command is an Action by its group but... if Mother is *just* wrapping the Cobra command tree, then the Action cannot have an `.Update()` or a `.View()` subroutine. We have identified the Action to `.Run()`, but we cannot associate new subroutines to it that can drive Bubble Tea while it is active.

With Type Embedding, an Action struct could embed cobra.Command and implement `.Update()` and `.View()` (basically: `class Action extends cobra.Command implements tea.Model` in OOP parlance). That way, it has all the subroutines Cobra will invoke in non-interactive mode and the two we need when driving Bubble Tea.

Solved, right? Not quite. The relationship must be bi-directional, which is not feasible.

Clock this signature `.AddCommand(cmds ...*cobra.Command)`. To get commands into its tree, we need to supply a cobra.Command *struct*. Due to the way Go's quasi-inheritance works, we cannot masquerade our Action 'super' type as its 'base'. We can supply cobra with a pointer to the embedded type. ex: 

```go
a := &action.Action{Command: cobra.Command{}}

root.AddCommand(a.Command)
```

This, however, will dispose of our super wrapper `a` as soon as it falls out of scope, meaning we must maintain a secondary data structure to hold onto the pointer to `a`.

We have two options:

1) Wrap all entry points into the cobra tree. This functionally maintains a second tree, which we want to avoid, especially given that Navs require no additional data (you'll notice they are aliased to cobra.Commands internally). Whether we maintain an entirely separate tree or "just" wrap the entrypoints, the end result will serve to further decouple Cobra from Bubble Tea, which, again, is less than ideal.

2) Maintain a data structure of Actions within Mother so we can look up the Update and View subroutines of the requested Action when it is relevant.

While it certainly isn't ideal, I believe 2 is the more maintainable option. We still, functionally, have to register the Actions twice at start up: once to cobra and once to Mother, despite mother containing cobra. Not great, but it means Bubble Tea/interactive mode can function entirely off Cobra's navigation and Cobra can operate entirely as normal. The only adaptation takes place in interactive mode, when an action is invoked; Mother uses the action cobra.Command to fetch the `.Update()` and `.View()` functions that should supplant her standard model.

*If you can figure a better adaption pattern, I am all ears.*