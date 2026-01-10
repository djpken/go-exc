# Requirements

## REQ-01: 修改 Module Name
- 將專案的 module name 從 `github.com/djpken/go-okex` 修改為 `github.com/djpken/go-exc`
- 更新 go.mod 文件
- 全域替換所有 Go 文件中的 import 路徑

## REQ-02: 升級 WebSocket 依賴版本
- 將 `github.com/gorilla/websocket` 從 v1.4.2 升級到 v1.5.3
- 更新 go.mod 和 go.sum
- 確保代碼相容性

## REQ-03: 升級 Go 版本
- 將專案的 Go 版本從 1.21 升級到最新穩定版本
- 更新 go.mod 中的 go 版本聲明
- 注意：Prompt 中提到 1.25.1，但此版本不存在，需確認正確版本（建議 1.23.x）

## REQ-04: 執行專案重構計劃
- 按照 REFACTORING_PLAN.md 進行完整的專案重構
- 將 go-okex 專案重構為支援多交易所的 go-exc 專案

### REQ-04.1: 階段1 - 基礎重構與重命名
- 創建新目錄結構
- 移動現有代碼到 exchanges/okex/ 目錄
- 更新包內 import 路徑

### REQ-04.2: 階段2 - 核心接口實現
- 創建核心接口文件（exc.go, config.go, errors.go, factory.go）
- 創建通用數據模型（types/包）

### REQ-04.3: 階段3 - OKEx 適配器實現
- 實現 OKEx Exchange 接口
- 創建類型轉換器
- 實現 REST 和 WebSocket 適配器

### REQ-04.4: 階段4 - 示例和文檔
- 創建使用示例代碼
- 更新專案文檔（README.md, MIGRATION.md, ARCHITECTURE.md）

### REQ-04.5: 階段5 - 測試與驗證
- 執行編譯驗證
- 創建並運行測試
- 確保所有功能正常運作
