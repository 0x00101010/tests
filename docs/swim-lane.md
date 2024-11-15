```mermaid
sequenceDiagram
    participant Customer
    participant Store
    participant Payment
    participant Inventory

    Customer->>Store: Browse Products
    Store->>Inventory: Check Stock
    Inventory-->>Store: Stock Status
    Store-->>Customer: Show Availability
    Customer->>Store: Add to Cart
    Customer->>Payment: Checkout
    Payment->>Store: Process Payment
    Store->>Inventory: Update Stock
    Store-->>Customer: Confirm Order
```