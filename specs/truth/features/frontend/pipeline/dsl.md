# Frontend Pipeline Module DSL

> **模組 DSL (Frontend: Pipeline View)**  
> **Truth Root**: `specs/truth/features/frontend/pipeline/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^管線看板各階段合計應顯示：$` | Then | DataTable: `[階段, 筆數]` | 無 | 讀取各看板欄位 header 的 count badge，斷言筆數符合 DataTable。 |
| `^管線看板應僅顯示 "([^"]*)" 負責之卡片$` | Then | `ownerName: string` | 無 | 斷言畫面上所有 Opportunity Card 的負責人標籤皆為 `{ownerName}`。 |
| `^不應存在 "([^"]*)" 負責之卡片$` | Then | `unauthorizedOwner: string` | 無 | 斷言畫面上沒有任何卡片標註負責人為 `{unauthorizedOwner}`。 |
