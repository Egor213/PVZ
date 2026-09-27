from pathlib import Path

with open("all_py.txt", "w", encoding="utf-8") as out:
    for file in Path(".").rglob("*.go"):
        if "test" not in str(file) and "repomocks" not in str(file) and "servicemocks" not in str(file):
            out.write(f"\n\n# ===== FILE: {file} =====\n\n")
            out.write(file.read_text(encoding="utf-8"))
