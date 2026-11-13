# B2B Lien Waiver & Payment Engine

A vertical SaaS backend built to automate legal compliance and escrow payments for commercial construction. This system dynamically generates state-specific statutory legal documents and routes B2B payments via Stripe Connect FBO (For-Benefit-Of) accounts, eliminating money transmitter liability.

## Architecture
* **Frontend:** React
* **Backend:** Go
* **Database:** Postgres
* **Integrations:** Stripe Connect, Twilio

## Core Technical Achievements
* **Cryptographic ESIGN Audit Trail:** Captures IP, device metadata, and SMS 2FA tokens to hash against generated PDFs for legal non-repudiation.
* **Idempotent Webhook Processing:** Ensures exact-once execution of critical financial settlement events from Stripe.
* **Dynamic PDF Mapping:** Maps relational database fields to strict X/Y bounding boxes on state-mandated statutory forms.

## How to Run Locally
1. `git clone ...`
2. `npm install`
3. Create a `.env` file (see `.env.example`)
4. `npm run dev`