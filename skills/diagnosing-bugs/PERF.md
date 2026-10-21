# Performance Diagnosis & Benchmark Vetting

Read this when a bug is "slow", a change claims a speedup or regression, or any test or report depends on a performance number. Answer each question with evidence from a run, not a guess about the code.

## Before You Run Anything

- Write down the claim you expect to make in the words you would ship ("export is 30% faster at p50 on the 60k-row dataset"). The questions below test that sentence.
- Read the measurement script: note what it times, counts, and ignores.
- Check machine load (`uptime`, core count). If it is busy and you cannot stop the noise, interleave the sides so both see it, and say so.

## Seven Questions

1. **Why not double?** Name the limiter (CPU per process, a runtime profiler, I/O wait, syscall counts) and map the hot spot to source. Profile in a run you do not report, because profilers slow the work. If the load generator saturates first, you measured the load generator.
2. **Was it tuned?** Run every side the way production runs it: release builds, production flags, pools, batching, cache temperature, versions, and data. A limiter that is a setting (a commit per row, a debug build, a missing index) means that side is untuned; tune it and measure again before picking a winner.
3. **Did it break limits?** Do the arithmetic against disk and network bandwidth and core count. Removing a piece that takes 10% of the run can make it at most about 11% faster. A result past a limit measured a cache, a no-op, or a bug.
4. **Did it error?** Count failures and non-success responses, and check outputs are correct, not just present. Rejections are fast; timeouts and retries are slow.
5. **Does it reproduce?** Run each side at least 5 times, alternating A, B, A, B, and report the median and range. A gap smaller than run-to-run variation is no measurable difference.
6. **Does it matter?** Next to any micro result, measure the end-to-end path a user waits on and report the micro result as a share of the whole.
7. **Did it even happen?** Confirm the work ran inside the timed region: the request reached the server, rows were written, the result was used. Lazy code and timeouts produce numbers for work that never happened.

For a quick ballpark the user asked for, one run is enough: still check questions 4 and 7, and say it is one run.

## Fix Order (Cheapest First)

Try these in order and stop when one meets the target:

1. Don't do it: stop work whose result nothing uses.
2. Do it, but don't do it again (cache or memoize).
3. Do it less.
4. Do it later.
5. Do it when no one is waiting.
6. Do it concurrently.
7. Do it cheaper.

Verify each attempt before trying the next.

## Report

Lead with the verdict: faster, slower, no measurable difference, or **inconclusive**. Give the number with its unit, the run count, the range, and the limiter ("p50 41 ms -> 33 ms, median of 7 runs per side, range 32–35 ms, bound by JSON parsing on one core"). The verdict is inconclusive when you cannot name the limiter, a side ran untuned, or questions 4 or 7 went unchecked; name the gap. Inconclusive is not a pass.
