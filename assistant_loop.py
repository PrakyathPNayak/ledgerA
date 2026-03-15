#!/usr/bin/env python3
"""Ralph Wiggum assistant loop — reads TODO.md, runs backpressure gates, and manages task flow."""

import os
import sys
import subprocess
from pathlib import Path

COLORS = {
    'green':  '\033[92m', 'yellow': '\033[93m',
    'red':    '\033[91m', 'reset':  '\033[0m',
    'bold':   '\033[1m',  'cyan':   '\033[96m'
}


def c(color: str, text: str) -> str:
    """Wrap text in ANSI color codes."""
    return f"{COLORS[color]}{text}{COLORS['reset']}"

def run_gate(name: str, cmd: str, cwd: str | None = None) -> bool:
    """Run a single backpressure gate command and return True if it passes."""
    print(f"  {c('cyan', '→')} {name:<40}", end="", flush=True)
    home = str(Path.home())
    env = os.environ.copy()
    env["GOPATH"] = f"{home}/go"
    env["GOMODCACHE"] = f"{home}/go/pkg/mod"
    env["GOCACHE"] = f"{home}/.cache/go-build"
    env["PATH"] = f"{home}/go/bin:" + env.get("PATH", "")
    Path(env["GOMODCACHE"]).mkdir(parents=True, exist_ok=True)
    Path(env["GOCACHE"]).mkdir(parents=True, exist_ok=True)
    try:
        result = subprocess.run(
            cmd, shell=True, cwd=cwd, capture_output=True, text=True, env=env
        )
        if result.returncode == 0:
            print(c('green', ' PASS'))
            return True
        print(c('red', ' FAIL'))
        for line in (result.stdout + result.stderr).strip().splitlines()[:25]:
            print(f"       {c('red', line)}")
        return False
    except Exception as exc:
        print(c('red', f' ERROR: {exc}'))
        return False


def check_backpressure() -> bool:
    """Run all applicable backpressure gates. Returns True if all pass."""
    all_pass = True
    print(f"\n{c('bold', 'Backpressure Gates:')}")

    if Path("go.mod").exists():
        if not run_gate("go build ./...", "go build ./..."):
            all_pass = False
        if not run_gate("go vet ./...", "go vet ./..."):
            all_pass = False
        if not run_gate("go test ./... -race -count=1", "go test ./... -race -count=1"):
            all_pass = False
        if not run_gate("golangci-lint",
                        "golangci-lint run ./cmd/... ./internal/... ./pkg/..."):
            all_pass = False

    if Path("frontend/package.json").exists():
        if not run_gate("tsc --noEmit", "npx tsc --noEmit", cwd="frontend"):
            all_pass = False
        if not run_gate("eslint src/", "npx eslint src/", cwd="frontend"):
            all_pass = False
        if not run_gate("vite build", "npm run build", cwd="frontend"):
            all_pass = False

    print()
    return all_pass


def main() -> None:
    """Main entry point for the assistant loop."""
    todo_file = Path("TODO.md")
    task_file = Path("RALPH_TASK.md")

    print(f"\n{c('bold', '═══════════════════════════════════════')}")
    print(f"{c('bold', '   RALPH WIGGUM TASK MANAGER   ')}")
    print(f"{c('bold', '═══════════════════════════════════════')}\n")
    print(f"  {c('cyan', 'Frontend')}: http://localhost:5173")
    print(f"  {c('cyan', 'API')}: http://localhost:8080/api/v1/health\n")

    if not todo_file.exists():
        print(c('yellow', 'WARNING: TODO.md not found — gates-only mode.'))
        gates_passed = check_backpressure()
        if not gates_passed:
            print(c('red', '✗ Backpressure gates failed. Fix errors before proceeding.'))
            sys.exit(1)
        try:
            user_input = input(f"{c('bold', 'ENTER to continue | type instruction: ')}").strip()
            if user_input:
                task_file.write_text(f"Override Instruction:\n{user_input}\n")
            else:
                task_file.write_text("Active Task:\nNo TODO.md task available — gates-only mode.\n")
        except KeyboardInterrupt:
            print("\nExiting.")
        sys.exit(0)

    raw_lines = todo_file.read_text().splitlines()
    completed: list[str] = []
    pending: list[str] = []

    for line in raw_lines:
        stripped = line.strip()
        if stripped.startswith("- [x]"):
            completed.append(stripped)
        elif stripped.startswith("- [ ]"):
            pending.append(stripped)

    if completed:
        print(f"{c('green', '✓ Last completed:')} {completed[-1][6:]}")
    else:
        print(c('yellow', '  No completed tasks yet.'))

    if not pending:
        print(f"\n{c('bold', c('green', '🎉  ALL TASKS COMPLETE  🎉'))}\n")
        sys.exit(0)

    next_task = pending[0]
    print(f"{c('yellow', '→ Next task:  ')} {next_task[6:]}\n")

    gates_passed = check_backpressure()
    if not gates_passed:
        print(c('red', '✗ Backpressure gates failed. Fix all errors before proceeding.\n'))
        sys.exit(1)

    print(c('green', '✓ All gates passed.\n'))

    try:
        user_input = input(f"{c('bold', 'ENTER to continue | type instruction: ')}").strip()
        if user_input:
            task_file.write_text(f"Override Instruction:\n{user_input}\n")
        else:
            task_file.write_text(f"Active Task:\n{next_task}\n")
    except KeyboardInterrupt:
        print("\nExiting.")
        sys.exit(1)

    sys.exit(0)


if __name__ == "__main__":
    main()
