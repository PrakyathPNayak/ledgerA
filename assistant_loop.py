#!/usr/bin/env python3
import os
import sys, subprocess, json, re
from pathlib import Path

COLORS = {
    'green':  '\033[92m', 'yellow': '\033[93m',
    'red':    '\033[91m', 'reset':  '\033[0m',
    'bold':   '\033[1m',  'cyan':   '\033[96m'
}

def run_gate(name, cmd, cwd=None):
    try:
        print(f"Running gate: {name}...", end=" ", flush=True)
        env = os.environ.copy()
        home = str(Path.home())
        env["GOPATH"] = f"{home}/go"
        env["GOMODCACHE"] = f"{home}/go/pkg/mod"
        env["GOCACHE"] = f"{home}/.cache/go-build"
        env["PATH"] = f"{home}/go/bin:" + env.get("PATH", "")
        Path(env["GOMODCACHE"]).mkdir(parents=True, exist_ok=True)
        Path(env["GOCACHE"]).mkdir(parents=True, exist_ok=True)

        result = subprocess.run(cmd, shell=True, cwd=cwd, capture_output=True, text=True, env=env)
        if result.returncode == 0:
            print(f"{COLORS['green']}PASS{COLORS['reset']}")
            return True
        else:
            print(f"{COLORS['red']}FAIL{COLORS['reset']}")
            print(f"{COLORS['red']}{result.stdout}{result.stderr}{COLORS['reset']}")
            return False
    except Exception as e:
        print(f"{COLORS['red']}ERROR: {e}{COLORS['reset']}")
        return False

def check_backpressure():
    gates_passed = True
    
    # Check if backend exists
    if Path("go.mod").exists():
        if not run_gate("go build", "go build ./..."): gates_passed = False
        if not run_gate("go vet", "go vet ./..."): gates_passed = False
        if not run_gate("go test", "go test ./... -race -count=1"): gates_passed = False
        if not run_gate("golangci-lint", "golangci-lint run"): gates_passed = False

    # Check if frontend exists
    if Path("frontend/package.json").exists():
        if not run_gate("tsc", "npx tsc --noEmit", cwd="frontend"): gates_passed = False
        if not run_gate("eslint", "npx eslint src/", cwd="frontend"): gates_passed = False
        if not run_gate("vite build", "npm run build", cwd="frontend"): gates_passed = False

    return gates_passed

def main():
    todo_file = Path("TODO.md")
    if not todo_file.exists():
        print(f"{COLORS['yellow']}TODO.md not found. Running gates-only mode.{COLORS['reset']}")
        print(f"{COLORS['cyan']}Link:{COLORS['reset']} Frontend http://localhost:5173")
        print(f"{COLORS['cyan']}Link:{COLORS['reset']} API health http://localhost:8080/api/v1/health")

        if Path("go.mod").exists() or Path("frontend/package.json").exists():
            if not check_backpressure():
                print(f"{COLORS['red']}Backpressure gates failed. Fix errors before proceeding.{COLORS['reset']}")
                sys.exit(1)

        try:
            user_input = input("ENTER to continue | type instruction: ")
            task_file = Path("RALPH_TASK.md")
            if user_input.strip() == "":
                task_file.write_text("Active Task:\nNo TODO.md task available")
            else:
                task_file.write_text(f"Override Instruction:\n{user_input}")
            sys.exit(0)
        except KeyboardInterrupt:
            print("\nExiting...")
            sys.exit(1)

    todos = todo_file.read_text().splitlines()
    completed = []
    pending = []
    
    for line in todos:
        if line.strip().startswith("- [x]"):
            completed.append(line)
        elif line.strip().startswith("- [ ]"):
            pending.append(line)
            
    if completed:
        print(f"{COLORS['green']}Last completed:{COLORS['reset']} {completed[-1]}")
        
    if not pending:
        print(f"{COLORS['bold']}{COLORS['green']}ALL TASKS COMPLETE{COLORS['reset']}")
        sys.exit(0)
        
    next_task = pending[0]
    print(f"{COLORS['yellow']}Next task:{COLORS['reset']} {next_task}")
    print(f"{COLORS['cyan']}Link:{COLORS['reset']} Frontend http://localhost:5173")
    print(f"{COLORS['cyan']}Link:{COLORS['reset']} API health http://localhost:8080/api/v1/health")
    
    # Run backpressure gates only if there's actually code (skip for early tasks)
    if Path("go.mod").exists() or Path("frontend/package.json").exists():
        if not check_backpressure():
            print(f"{COLORS['red']}Backpressure gates failed. Fix errors before proceeding.{COLORS['reset']}")
            sys.exit(1)
            
    try:
        user_input = input(f"ENTER to continue | type instruction: ")
        
        task_file = Path("RALPH_TASK.md")
        if user_input.strip() == "":
            task_file.write_text(f"Active Task:\n{next_task}")
        else:
            task_file.write_text(f"Override Instruction:\n{user_input}")
            
        sys.exit(0)
    except KeyboardInterrupt:
        print("\nExiting...")
        sys.exit(1)

if __name__ == "__main__":
    main()
