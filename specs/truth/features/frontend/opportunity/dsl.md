# Frontend Opportunity Module DSL

> **模組 DSL (Frontend: Opportunity Stepper & RBAC)**  
> **Truth Root**: `specs/truth/features/frontend/opportunity/dsl.md`  
> **維護者**: `/dsl-refine`

| 句型 (Gherkin Pattern) | 類型 | 參數 / DataTable | 預設值 | StepDef 實作語意 |
|---|---|---|---|---|
| `^嘗試點擊推進階段按鈕 "([^"]*)"$` | When | `btnText: string` | 無 | 點擊標籤包含 `{btnText}` 的按鈕。 |
| `^畫面步進器維持在 "([^"]*)"$` | Then | `currentStage: string` | 無 | 檢驗 Stage Stepper 目前 Active/Current 的步進節點標籤為 `{currentStage}`。 |
