# Payments

There are ways to receive payments, through payment processors and MoR.

## Payment processor (Stripe)
Its main function is to move money from the client's account to yours. However, this makes you the **Seller of Record**, which means that you are legally responsible for calculating, collecting and reporting taxes (VAT/Sales Tax), issuing legal invoices and managing the tax file.

## Merchant of Record (MoR)
It is an entity that assumes complete fiscal responsibility. The MoR becomes the legal Seller of Record before the authorities. This frees you from the administrative burden of calculating international taxes, issuing invoices and managing returns or chargebacks.

## Feature comparison

| Feature | Payment Processor (e.g. Stripe) | Merchant of Record (e.g. Polar, Paddle) |
|---|---|---|
| Fiscal responsibility | The owner of the SaaS (Seller of Record) | The MoR takes responsibility |
| Tax management | Manual or with external tools (Stripe Tax) | Automatic and managed by the MoR |
| Legal billing | User Responsibility | MoR issues legal invoices |
| Commissions | Minor (~2.9% + 0.30) | Seniors (4% - 8%) |
| Chargebacks | User risk and management | The MoR usually manages them |

MoR has higher commissions, the cost is justified by delegating the "fiscal plumbing" and avoiding legal risks.
As an independent developer I will be working with MoR.