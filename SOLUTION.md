# Solution Description

## Algorithm (Critical Path Method)

The job and its tasks form a **directed acyclic graph (DAG)** where each node is a task and edges represent dependencies.

I use **Critical Path Method (CPM)** to compute the minimum completion time:

1. **Topological sort** of all tasks (Kahn's algorithm).
2. **Forward pass** to compute, for each task:
   - `EST(task) = max( EFT(dep) for all dependencies )` (or `0` if there is no dependency),
   - `EFT(task) = EST(task) + duration(task)`.
3. The **minimum completion time** of the job is `max( EFT(task) )` over all tasks.
4. The **critical path** is found by tracing backwards from the task with the largest `EFT`, always choosing a dependency whose `EFT` equals the current task's `EST`.

This runs in **O(V + E)** time, where `V` is the number of tasks and `E` is the number of dependency edges.

Assumption: each task uses a single worker, but independent tasks can run in parallel once their dependencies are finished.
