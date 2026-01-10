# Main entry point

## Principle

### 1. 優先使用 MCP Server 了解專案

**步驟 A：初步探索專案結構**
- `mcp__jetbrains__get_project_modules` - 獲取專案模組資訊
- `mcp__jetbrains__list_directory_tree` - 查看目錄結構
- `mcp__jetbrains__get_project_dependencies` - 檢查專案依賴
- `mcp__jetbrains__get_run_configurations` - 查看可用的運行配置
- `mcp__jetbrains__get_all_open_file_paths` - 查看當前打開的檔案

**步驟 B：深入程式碼分析**（針對特定功能或模組）
- `mcp__jetbrains__search_in_files_by_text` - 搜尋相關程式碼
- `mcp__jetbrains__read_file` - 讀取特定檔案內容
- `mcp__jetbrains__get_symbol_references` - 查找符號使用位置
- `mcp__jetbrains__get_type_hierarchy` - 了解類別繼承關係

**步驟 C：開發前檢查**
- `mcp__jetbrains__get_file_problems` - 檢查現有錯誤和警告

**重點**：確保完全了解專案架構後再進行開發工作

---

### 2. 開發流程

**步驟 A：程式碼修改**
- 使用 `mcp__jetbrains__apply_diff` 應用程式碼變更
- 修改時確保註解使用英文

**步驟 B：即時驗證**
- 使用 `mcp__jetbrains__get_file_problems` 檢查新增的錯誤
- 必要時使用 `go build` 檢查編譯問題

**步驟 C：測試與驗證**
- 使用 `mcp__jetbrains__run_configuration` 執行測試或運行程式
- 根據 Task 中 _Test 欄位進行對應驗證

---

### 3. 依序完成文件中的 [Prompt](#prompt)，整理到 [Requirements](/.code/Requirements.md) 中

**範例格式：**
```markdown
## REQ-16: ActionButton 添加更新資訊按鈕
- 在 ActionButtons 組件添加「更新資訊」按鈕
- 點擊後顯示更新資訊對話框
```

---

### 4. 將 [Requirements](/.code/Requirements.md) 拆解成子 task 放入 [Tasks](/.code/Tasks.md)

**子 task 格式規範：**
```markdown
- [x] 1. 初始化專案結構和基礎配置
  - 使用 Create React App with TypeScript 模板建立專案
  - 安裝核心依賴：Redux Toolkit, @dnd-kit/core, Ant Design, fast-xml-parser
  - 配置 TypeScript、ESLint、Prettier
  - 建立基本的資料夾結構（components, store, services, types, utils）
  - _Requirement: REQ-1
  - _Status: Done
  - _Test: Manual
```

**欄位說明：**
- `_Requirement`: 對應 Requirements.md 中的需求編號（如 REQ-1, REQ-16）
- `_Status`: 任務狀態
    - `Todo` - 待開始
    - `In Progress` - 進行中
    - `Done` - 已完成
    - `Blocked` - 被阻擋（需註明原因）
- `_Test`: 驗證方式
    - `Manual` - 手動測試
    - `Unit Test` - 單元測試
    - `Integration Test` - 整合測試
    - `E2E Test` - 端對端測試

---

### 5. 文件管理策略

#### Requirements.md 分層管理

當 Requirements.md 文件過長（超過 50 個需求項目）時，採用模組化拆分：

```
/.code/
├── Requirements.md           # 主索引文件
├── requirements/
│   ├── auth.md              # 認證相關需求
│   ├── ui-components.md     # UI 組件需求
│   ├── data-management.md   # 數據管理需求
│   └── integration.md       # 整合功能需求
```

**Requirements.md 主文件格式：**
```markdown
# Requirements Index

## 認證模組 (REQ-001 ~ REQ-020)
詳見：[auth.md](./requirements/auth.md)
- 登入功能
- 權限管理
- 用戶管理

## UI 組件 (REQ-021 ~ REQ-050)
詳見：[ui-components.md](./requirements/ui-components.md)
- ActionButton 組件
- Dialog 組件
```

#### Tasks.md 分層管理

當 Tasks.md 文件過長（超過 30-40 個任務）時，採用歸檔策略：

```
/.code/
├── Tasks.md                 # 當前 Sprint/階段任務
├── tasks/
│   ├── archive/
│   │   ├── 2024-Q1.md      # 已完成任務歸檔
│   │   └── 2024-Q2.md
│   └── backlog.md          # 待排程任務
```

**Tasks.md 主文件格式：**（僅保留當前 Sprint）
```markdown
# Current Sprint Tasks (2024-W02)

## 進行中 (In Progress)
- [ ] 15. 實作用戶認證 API
  - _Requirement: REQ-001, REQ-002
  - _Status: In Progress
  - _Test: Unit Test

## 待開始 (Todo)
- [ ] 16. 建立登入頁面 UI
  - _Requirement: REQ-003
  - _Status: Todo
  - _Test: Manual

---

## 歷史任務
- [2024 Q4 已完成任務](./tasks/archive/2024-Q4.md)
- [待排程任務列表](./tasks/backlog.md)
```

**歸檔規則：**
- 每完成 20-30 個任務，將已完成任務移至 `tasks/archive/`
- 按季度或月份歸檔（如 2024-Q1.md, 2024-Q2.md）
- 歸檔文件保留完整資訊供日後查詢

**跨文件引用規範：**
- 使用相對路徑引用需求：`[REQ-001](./requirements/auth.md#req-001)`
- Task 的 _Requirement 欄位使用完整編號：`REQ-001, REQ-002`
- 需求編號全局唯一，即使分散在不同文件中

---

### 6. 文件語言規範

- Codebase 中所有註解使用**英文**
- [Tasks](/.code/Tasks.md) 和 [Requirements](/.code/Requirements.md) 使用**中文**
- Commit message 建議使用英文（可選）

---

## Prompt
