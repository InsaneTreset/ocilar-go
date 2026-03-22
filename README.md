# Ocilar Go SDK

Official Go client for the [Ocilar API](https://ocilar.com) — CAPTCHA solving and document extraction.

## Install

```bash
go get github.com/InsaneTreset/ocilar-go
```

## Quick Start

```go
package main

import (
    "encoding/base64"
    "fmt"
    "os"

    ocilar "github.com/InsaneTreset/ocilar-go"
)

func main() {
    client := ocilar.NewClient("sk-YOUR_KEY")

    // Test connectivity
    hello, _ := client.Hello()
    fmt.Println(hello)

    // Solve SAT CAPTCHA
    imgBytes, _ := os.ReadFile("captcha.png")
    img := base64.StdEncoding.EncodeToString(imgBytes)

    result, _ := client.SolveSAT(img)
    fmt.Println(result.Text)      // "2VBF39"
    fmt.Println(result.LatencyMs) // 67
    fmt.Println(result.TaskID)    // "tsk_abc123"

    // Extract data from CSF
    docBytes, _ := os.ReadFile("csf.pdf")
    doc := base64.StdEncoding.EncodeToString(docBytes)

    extracted, _ := client.ExtractCSF(doc)
    fmt.Println(extracted.Data) // map[rfc:XAXX010101000 nombre:...]
}
```

## Available Methods

### CAPTCHA Solving
- `SolveSAT(imageBase64)` — SAT Mexico
- `SolveIMSS(imageBase64)` — IMSS Mexico
- `SolveImage(imageBase64)` — Generic image
- `SolveRecaptchaV2(siteKey, siteURL)` — reCAPTCHA v2
- `SolveRecaptchaV3(siteKey, siteURL, action)` — reCAPTCHA v3
- `SolveHCaptcha(siteKey, siteURL)` — hCaptcha
- `SolveCloudflare(siteURL)` — Cloudflare Turnstile
- `SolveAudio(audioBase64)` — Audio CAPTCHA

### Document AI
- `ExtractCSF(documentBase64)` — Constancia de Situacion Fiscal
- `ExtractINE(documentBase64)` — INE / Voter ID
- `ExtractCFDI(documentBase64)` — CFDI Invoice
- `ExtractCURP(documentBase64)` — CURP
- `ExtractDomicilio(documentBase64)` — Proof of Address
- `ExtractNomina(documentBase64)` — Payroll Receipt
- `ExtractGeneric(documentBase64)` — Generic OCR

### Utilities
- `Hello()` — Test API key
- `GetBalance()` — Account balance and usage

## Free Tier

Every account gets **1,000 free solves** per CAPTCHA type and **50 free document extractions** per type. No credit card required.

## Links

- [Dashboard](https://console.ocilar.com)
- [API Docs](https://api.ocilar.com/api/v1/docs)
- [Pricing](https://ocilar.com/#pricing)
