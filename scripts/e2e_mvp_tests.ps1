Param()

$ErrorActionPreference = 'Stop'

function Status([string]$label, [string]$url) {
    try {
        $code = & curl.exe -s -o NUL -w "%{http_code}" $url
        Write-Output ("$label $code")
    } catch {
        Write-Output ("$label ERR")
    }
}

$base = 'http://localhost:8080'
$v1 = "$base/api/v1"

# Quick status checks for dashboard/metrics/alerts/ws
Status 'dashboard_business' "$base/dashboard/business"
Status 'metrics' "$base/metrics"
Status 'alerts' "$base/alerts"
Status 'alerts_all' "$base/alerts?all=true"
Status 'ws_status' "$v1/ws/status"

# Financial endpoints
Status 'financial_prices' "$v1/financial/companies/1/prices"
Status 'financial_latest' "$v1/financial/companies/1/price/latest"
Status 'financial_indicators' "$v1/financial/companies/1/indicators"
Status 'market_latest' "$v1/financial/market"
Status 'market_history' "$v1/financial/market/history?start_date=2025-08-01&end_date=2025-08-10&limit=5"

# Analytics endpoints
Status 'analytics_trends' "$v1/analytics/companies/1/esg-trends?days=5"
Status 'analytics_sector_comp' "$v1/analytics/sectors/comparisons"
Status 'analytics_fin_comp' "$v1/analytics/financial/comparisons?limit=5"
Status 'analytics_top_esg' "$v1/analytics/top-performers/esg_score?limit=5"
Status 'analytics_corr' "$v1/analytics/correlation/esg-financial"
Status 'analytics_summary' "$v1/analytics/summary"

# ESG CRUD flow (create -> get -> list by company -> update -> delete)
try {
    $esgCreateBody = @{
        company_id = 1
        environmental_score = 70.5
        social_score = 71.2
        governance_score = 69.8
        overall_score = 70.5
        score_date = '2025-08-11T00:00:00Z'
        data_source = 'TEST'
    } | ConvertTo-Json
    $createResp = Invoke-RestMethod -Uri "$v1/esg/scores" -Method Post -ContentType 'application/json' -Body $esgCreateBody
    $esgId = $createResp.data.id
    $byCompany = Invoke-RestMethod -Uri "$v1/esg/companies/1/scores?limit=2"
    $esgGet = Invoke-RestMethod -Uri "$v1/esg/scores/$esgId"
    $esgUpdateBody = @{
        environmental_score = 75.0
        social_score = 74.0
        governance_score = 70.0
        overall_score = 73.0
        score_date = '2025-08-12T00:00:00Z'
        data_source = 'TEST2'
    } | ConvertTo-Json
    $esgUpd = Invoke-RestMethod -Uri "$v1/esg/scores/$esgId" -Method Put -ContentType 'application/json' -Body $esgUpdateBody
    $delCode = & curl.exe -s -o NUL -w "%{http_code}" -X DELETE "$v1/esg/scores/$esgId"
    Write-Output ("ESG_OK id=$esgId del=$delCode by_company_count=" + ($byCompany.scores | Measure-Object | Select-Object -ExpandProperty Count))
} catch {
    Write-Output ("ESG_ERR " + $_.Exception.Message)
}


