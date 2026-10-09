
Create a doc on how to get started with paddle sonbox and production
- Success url
- webhooks 

Create section explaining how to add products to paddle and link it to our plans or products

Sandbox:
- create account sandox
    - https://sandbox-vendors.paddle.com/
- create products and prices
- Create Api key:
    - My account > Settings > Authentication > New API Key 
- Configure default link: Checkout>Checkout configuration>Default payment link
    - Paddle requires a public link, not a local one; you can use a Cloudflare Tunnel.
      go tool task run:tunnel
    - cloudflared_url
- Configure webhook:
    - Events > notifications > New destination
    - URL: cloudflared_url/api/v1/paddle/webhook
    - Description: some desc
    - Select "All" for events for the sake of simplicity.

Production
- set APP_ENV=prod
- create live account
    - https://vendors.paddle.com/