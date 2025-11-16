# OpenTDF-MCP Encryption Bug Fix - Submodule Changes

## Summary

The encryption bug fix was made in the `opentdf-mcp` submodule, which is a separate Git repository (`https://github.com/snpm/opentdf-mcp`). The changes are currently only committed locally and not pushed to any remote repository.

## The Problem

The main repository (`nibsbin/opentdf-memo-demo`) references the submodule at commit `d43e4f3`, but this commit only exists locally. When someone clones your repository and initializes submodules, they'll get the original commit `c2e30bd` from the upstream repository, not your fixed version.

## Files with Changes

The patch includes changes to these files in the `opentdf-mcp` submodule:
1. **BUG_FIX.md** (new file) - Documentation of the bug and fix
2. **mcp-server/main.go** - The actual code fix

## Viewing the Changes

The complete patch file is available at:
- `0001-Fix-encryption-bug-read-file-contents-instead-of-enc.patch`

You can review it with:
```bash
cat 0001-Fix-encryption-bug-read-file-contents-instead-of-enc.patch
```

## Options to Make Changes Accessible

### Option 1: Fork the opentdf-mcp Repository (Recommended)

1. Fork `https://github.com/snpm/opentdf-mcp` to your GitHub account
2. Update the submodule to point to your fork:
   ```bash
   cd /home/user/opentdf-memo-demo
   git config -f .gitmodules submodule.opentdf-mcp.url https://github.com/YOUR_USERNAME/opentdf-mcp
   git submodule sync
   ```
3. Push the submodule changes to your fork:
   ```bash
   cd opentdf-mcp
   git push origin HEAD:refs/heads/fix-encryption-bug
   ```
4. Update main repo to reference your fork and push

### Option 2: Apply Patch to Your Fork

If you already have a fork:
1. Clone your fork of opentdf-mcp
2. Apply the patch:
   ```bash
   cd your-opentdf-mcp-fork
   git am < /path/to/0001-Fix-encryption-bug-read-file-contents-instead-of-enc.patch
   git push origin main  # or feature branch
   ```
3. Update the submodule reference in the main repo

### Option 3: Create PR to Upstream

1. Fork the upstream repository
2. Apply the patch to your fork
3. Create a pull request to `https://github.com/snpm/opentdf-mcp`
4. Once merged, update your submodule reference

### Option 4: Include Changes Inline (Not Recommended)

Remove the submodule and include the opentdf-mcp code directly in your repository. This loses the benefits of using a submodule but makes the code directly accessible.

## Summary of Changes

### Code Changes (`mcp-server/main.go`)

1. **Modified `EncryptToolInput` struct** (lines 22-28)
   - Added `Input` field for file paths
   - Made `Data` field optional
   - Both fields are mutually exclusive

2. **Updated `MCPEncrypt` function** (lines 96-116)
   - Added validation to ensure exactly one parameter is provided
   - Implemented file reading logic when `input` is specified
   - Reads file contents using `os.ReadFile()` instead of encrypting the path string

3. **Updated tool description** (line 327)
   - Clarifies usage: "Specify either 'input' (file path) or 'data' (literal text)"

### Impact

This fix ensures that when users call:
```go
mcp__opentdf-mcp__encrypt(input: "/path/to/file.txt", format: "nano")
```

The function reads and encrypts the **file contents**, not the string "/path/to/file.txt".

## Current State

- Main repo commits:
  - `61ccc33` - Documentation updates (CLAUDE.md, memo-buddy.chatmode.md)
  - `5b6eaab` - Submodule reference update
- Submodule commit: `d43e4f3` (local only, not pushed anywhere)
- Branch: `claude/fix-opentdf-encryption-bug-01UAuvzCRynYaNuCjUdNeWT5`

## Recommended Next Steps

1. Decide how you want to handle the submodule (fork, PR upstream, or inline)
2. If forking: Create fork and push the changes
3. Update `.gitmodules` to point to your fork (if applicable)
4. Rebuild the `opentdf-mcp-server` binary with the fixes
5. Test the encryption with both file paths and literal data
