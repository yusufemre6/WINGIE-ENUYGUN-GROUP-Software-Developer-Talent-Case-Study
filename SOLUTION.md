# Solution Description

## Terms

- **DAG** — Directed Acyclic Graph: a graph with directed edges and no cycles (no path from a task back to itself).
- **CPM** — Critical Path Method: the algorithm used to find the minimum completion time.
- **EST** — Earliest Start Time: the earliest time a task can start (after all its dependencies finish).
- **EFT** — Earliest Finish Time: the earliest time a task can finish; `EFT = EST + duration`.
- **V** — number of tasks (vertices in the graph).
- **E** — number of dependency edges.

## Algorithm (Critical Path Method)

The job and its tasks form a **DAG** where each node is a task and edges represent dependencies.

1. **Topological sort** of all tasks (Kahn's algorithm) so that whenever task X depends on task Y, Y appears before X in the order.
2. **Forward pass**: process tasks in that order. For each task:
   - `EST = 0` if it has no dependencies, otherwise `EST = max(EFT of each dependency)`.
   - `EFT = EST + duration`.
3. **Minimum completion time** = `max(EFT)` over all tasks.
4. **Critical path**: from the task with the largest EFT, go backwards by always picking the dependency whose EFT equals the current task's EST.

Time complexity: **O(V + E)**.

Assumption: each task uses one worker; independent tasks can run in parallel once their dependencies are done.

---

## Worked Example (from the case study)

**Input**

| Task | Duration | Dependencies |
|------|----------|--------------|
| A    | 3        | —            |
| B    | 2        | —            |
| C    | 4        | —            |
| D    | 5        | A            |
| E    | 2        | B, C         |
| F    | 3        | D, E         |

So: A, B, C have no dependencies; D depends on A; E depends on B and C; F depends on D and E.

**Step 1 — Topological order**

Kahn's algorithm gives an order where every dependency is processed before the task that depends on it. One valid order is:

`A, B, C, D, E, F`

We will use this order for the forward pass.

**Step 2 — Forward pass (EST and EFT)**

| Task | Dependencies | EST | EFT |
|------|--------------|-----|-----|
| A    | none         | 0   | 0 + 3 = **3** |
| B    | none         | 0   | 0 + 2 = **2** |
| C    | none         | 0   | 0 + 4 = **4** |
| D    | A            | EFT(A) = 3 | 3 + 5 = **8** |
| E    | B, C         | max(EFT(B), EFT(C)) = max(2, 4) = **4** | 4 + 2 = **6** |
| F    | D, E         | max(EFT(D), EFT(E)) = max(8, 6) = **8** | 8 + 3 = **11** |

**Step 3 — Minimum completion time**

The job finishes when the last task finishes. The latest EFT is 11 (task F).

**Minimum completion time = 11** (unit times).

**Step 4 — Critical path**

Start from the task that finishes last: **F** (EFT = 11). F's EST is 8; the dependency that finishes at time 8 is **D**. So F is preceded by D. D's EST is 3; the dependency that finishes at time 3 is **A**. So the path is **A → D → F**.

**Critical path: A → D → F** (total duration 3 + 5 + 3 = 11). Any delay on this path increases the total completion time.
