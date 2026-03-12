# Mosaic - Tactical Profiler & Password Generator 🧩

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![Memory](https://img.shields.io/badge/Memory-Zero_Allocation-blueviolet?style=for-the-badge)
![Speed](https://img.shields.io/badge/Speed-0.15s_Generation-00ADD8?style=for-the-badge&logo=go)

**Mosaic** is a high-performance wordlist generator focused on OSINT and human profiling, heavily inspired by the legendary **![CUPP](https://github.com/Mebus/cupp)** tool. Developed in **Go**, it replaces blind brute-force generation with behavioral heuristics, processed through a highly concurrent architecture (![Goroutines](https://go.dev/tour/concurrency/1) and MPSC Channels).

## 🧠 The Math Behind the Engine (Why Mosaic is Superior)

Classic dictionary generators fail because they use a **Naive Cartesian Product** of all user inputs. If a target has a set of keywords $K$, and a set of dates $D$, a naive approach generates combinations like:
$$C_{naive} = K \times D \times K \dots$$
This generates a combinatorial explosion of "noise" (e.g., `15_Alice`, `07.Bob`), wasting CPU time during generation and network time during attacks.
 
### The Heuristic Modification (Guided Cartesian Product)
Mosaic implements a **3-Level Combinatorial Model** based on human psychology. We define the Base Words set $B$ (Names + Date Fragments) and the Common Suffixes set $S$ (e.g., `123`, `!`, years). The generation function creates the password set $P$ through the union of high-probability subsets:

$$P = (B \times S) \cup (B \times B) \cup (B \times B \times S)$$

This drastically prunes the search tree, ensuring that unrealistic combinations are never allocated in memory, focusing computational power on creating hyper-realistic targets like `Aj_15901990` (where $b_1 \in B$, $b_2 \in B$, and $s_1 \in S$).

## ⏱️ Asymptotic Analysis (Complexity)

### 1. Time Complexity
* **Base Heuristic Generation:** Let $|B|$ be the number of bases and $|S|$ the number of suffixes. The Level 3 loop runs in $\mathcal{O}(|B|^2 \cdot |S|)$. Since $|B|$ and $|S|$ are bounded by human profile inputs (usually $< 100$), this executes in milliseconds.
* **O(1) Deduplication:** We use a Go `map[string]struct{}`. Checking and inserting unique strings occurs in $\mathcal{O}(1)$ average case.
* **Leetspeak Mutation (Backtracking):** The most computationally heavy task. For a password of length $N$ with a maximum branching factor of leet options $L$ (e.g., 'a' -> 'a', '@', '4' means $L=3$), the complexity is $\mathcal{O}(L^N)$. Thanks to concurrent Workers, this load is distributed across $C$ CPU cores, resulting in a real processing time of $\mathcal{O}(\frac{L^N}{C})$.

### 2. Space Complexity
* Base memory usage is dominated by the deduplication map, taking $\mathcal{O}(U)$, where $U$ is the number of unique base strings generated.
* The Leetspeak algorithm was written as **Zero-Allocation**. We operate directly on a byte slice (`[]byte`), mutating indexes in-place and reverting them (backtracking). Local space complexity: $\mathcal{O}(N)$ for the recursive call stack, generating zero garbage for the Go GC.

## 🚀 Architecture and Performance
* **Worker Pool:** Mosaic automatically detects `runtime.NumCPU()` and spawns symmetric mutation Workers matching your physical/logical cores.
* **MPSC (Multi-Producer, Single-Consumer):** Multiple Workers process strings and inject them into the buffered `passChan`. A single consumer Goroutine handles the `bufio.Flush()` directly to the disk, eliminating Race Conditions and I/O blocking.

## 📊 Performance Benchmarks: Go vs Python

To ensure Mosaic provides a tactical advantage in real-world scenarios, we benchmarked it against the industry standard (CUPP - Python 3) using a strict 1:1 heuristic generation test. 

*Note: The test was conducted without Leetspeak mutations to ensure a fair baseline comparison of the core generation engine.* 

> [!NOTE]
> The tool used to mensure the time between both programs was ![time](https://man7.org/linux/man-pages/man1/time.1.html)

| Metric | CUPP Original (Python 3) | Mosaic (Go 1.20+) | Advantage |
| :--- | :--- | :--- | :--- |
| **Execution Time** | 0.42 Seconds | **0.15 Seconds** | **~2.8x Faster** |
| **Memory Peak (RAM)**| 30.4 MB | **25.0 MB** | **17% Less RAM** |
| **CPU Utilization** | 99% (GIL Limit - 1 Core) | **312% (Multi-Core)** | **Perfect Scaling** |
| **Workload (I/O Writes)**| 272 blocks written | **1,712 blocks written** | **Generated ~6.2x more permutations** |

### 🛠️ Why does Mosaic outperform?

1. **Breaking the GIL:** Python is limited by the Global Interpreter Lock (GIL), pinning execution to a single core. Mosaic uses Go's Goroutines to detect your `runtime.NumCPU()` and distribute mathematical loads across all available physical and logical cores simultaneously.
2. **Asynchronous I/O via MPSC:** Disk writing is the ultimate bottleneck. Mosaic implements a Multi-Producer Single-Consumer channel architecture. CPU workers push passwords to memory buffers in microseconds, while a single, dedicated Goroutine safely flushes the buffer to disk without race conditions.
3. **Zero-Allocation Mutations:** Deep string permutations stress the Garbage Collector. Mosaic performs in-place byte slice (`[]byte`) manipulation with backtracking algorithms, generating near-zero garbage.

## 🛠️ Roadmap (Future Generation Improvements)

While Mosaic is highly optimized, there are theoretical boundaries that can be evolved:
1. **Memory Management for Massive Dictionaries:** Currently, `map[string]struct{}` holds all permutations in RAM before sending them to Workers. For bulk keywords, a *Bloom Filter* could be implemented for streaming deduplication.
2. **Context-Aware Leetspeak:** The current backtracking replaces all occurrences blindly. A heuristic entropy mutation would limit substitutions only to the first or last vowel, drastically cutting the $\mathcal{O}(L^N)$ time for long passwords.
3. **Regional Pattern Injection (Keyboard Walks):** Humans frequently add "qwer" or "1q2w" as connectors. Adding keyboard walks to Level 2 joins would increase the hit probability against corporate targets.

## 🎮 Installation and Usage

**Fast Build (Optimized):**
```bash
go build -ldflags="-w -s" -o mosaic main.go
