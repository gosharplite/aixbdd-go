# Frontend Root Shared DSL

> **介面根共用 DSL (Frontend: Playwright E2E)**  
> **Truth Root**: `specs/truth/features/frontend/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^使用者在瀏覽器開啟 "([^"]*)" 頁面$` | Given | `pagePath: string` | 無 | Playwright `page.goto(pagePath)` 開啟指定頁面。 |
| `^使用者切換當前角色身分為 "([^"]*)"$` | When | `roleName: string` | 無 | 點選畫面上的角色切換器（Role Switcher）下拉或按鈕選取 `{roleName}`。 |
| `^畫面應顯示操作拒絕警告提示 "([^"]*)"$` | Then | `expectedAlert: string` | 無 | 驗證畫面上彈出 Toast、Modal 或 Alert 元素，其文字內容包含 `{expectedAlert}`。 |
| `^畫面標題應為 "([^"]*)"$` | Then | `title: string` | 無 | 斷言頁面標題或主要 h1 元素文字包含 `{title}`。 |
