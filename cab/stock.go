package main

/*
 Create a stock exchange system matching buy and sell orders efficiently.

 Requirements
 -> User can select a share, check for the price.
 -> User can buy the share for a price.
 -> User can sell the share for a price.

 Flows
 User checking the share (share_id)
 -> User selects a share
 -> Fetch the share details -> price, share_name, exchange

 User buying a share (share_id, price, operation= buy/sell, quantity )
	-> Create a order with the details with status = created
	-> Check for a corresponding order of sell/buy with that price.
		-> if available
		-> mark the order as partial/fully executed based on quantities of both the order.
		-> not available
			-> no corresponding order available, return the response with order is pending, will be executed in sometime.

Matching Algo
order data stored on per share_id so the fetch will be for only a particular share O(N), N = number of shares
 -> linear where check all the orders and see if it matches with the price
 -> use the sorted set to store the data in sorted order, lookup will be O(Logn) and insert/delete will be O(LogN)

 Entities
Order
- id
- share_id
- operationType (BUY, SELL)
- created_at
- quantities
- price

Share
- id
- availableQuantities
- price

PendingOrder
-> map[share_id][]Order(sorted set so that the insert, delete and search is O(LogN))

Services(Methods)
OrderService
- CreateOrder
- UpdateOrder

ShareService
- UpdateShare
- CreateShare
- DeleteShare

OrderPendingService
- CreateOrderPending
- FetchOrderForShareAndPrice
- DeleteOrderPending

NotificationService
-> Notify() //Notify the other User about the order execution

OrderService(Fields)
- OrderPendingService -> FetchOrderForShareAndPrice
- NotificationService
- OrderRepo
- ShareRepo

ShareService
- ShareRepo

OrderPendingService
- MatchingAlgo(Linear, Sorted Sets)
- OrderPendingRepo

NotificationService
- Notify(channel)

Design Patterns
-> Startegy Pattern for MatchingAlgo in OrderPendingService
-> Observer in NotificationService
-> Factory Pattern in creating the Matching Algo(Linear, Sorted Sets)

==========================================
SUGGESTIONS & IMPROVEMENTS
==========================================

MISSING ENTITIES:
1. User
   - id, name, email, balance
   - portfolio (holdings)

2. Trade (Executed Order)
   - id, buyOrderID, sellOrderID, share_id, quantity, price, executed_at

3. Portfolio/Holdings
   - user_id, share_id, quantity, averagePrice

4. Order Status (enum)
   - PENDING, PARTIAL, FILLED, CANCELLED

5. Order Type (enum)
   - LIMIT (specific price), MARKET (best available)

MISSING SERVICES:
1. MatchingService (separate from OrderPendingService)
   - FindMatch(order) -> *Order
   - MatchOrder(order) -> []Trade
   - Handle price-time priority

2. PortfolioService
   - GetPortfolio(userID)
   - UpdateHoldings(userID, shareID, quantity)
   - ValidateSell(userID, shareID, quantity)

3. PaymentService
   - ValidateBalance(userID, amount)
   - DeductBalance(userID, amount) // for buy
   - AddBalance(userID, amount) // for sell

IMPROVED FLOWS:
1. Buy Order Flow:
   -> Validate user balance
   -> Create order (status: PENDING)
   -> Try matching
   -> If match: Execute trade, update portfolio, notify
   -> If no match: Add to pending orders

2. Sell Order Flow:
   -> Validate user has shares
   -> Create order (status: PENDING)
   -> Try matching
   -> If match: Execute trade, update portfolio, notify
   -> If no match: Add to pending orders

3. Cancel Order Flow:
   -> Validate order belongs to user and is PENDING
   -> Remove from pending orders
   -> Update status to CANCELLED

4. Order Matching Flow:
   -> Check opposite order book (buy checks sell, sell checks buy)
   -> Match by price (best price first)
   -> Match by time (first come first serve)
   -> Handle partial fills
   -> Create Trade records
   -> Update both orders (FILLED/PARTIAL)
   -> Update share price (last traded)

IMPROVED ENTITIES:
Order:
- id, user_id, share_id
- operationType (BUY, SELL)
- orderType (LIMIT, MARKET)
- status (PENDING, PARTIAL, FILLED, CANCELLED)
- quantity, filledQuantity
- price
- created_at, updated_at

Share:
- id, name, symbol, exchange
- currentPrice (last traded price)
- totalShares (outstanding)

OrderBook:
- share_id
- buyOrders (sorted: price desc, time asc)
- sellOrders (sorted: price asc, time asc)

IMPROVED SERVICES:
OrderService (Orchestrator):
- CreateOrder(userID, shareID, operation, quantity, price, orderType)
- CancelOrder(userID, orderID)
- GetOrder(userID, orderID)
- GetUserOrders(userID, status)
- ExecuteTrade(buyOrder, sellOrder, quantity, price)

MatchingService:
- FindMatch(order) -> *Order
- MatchOrder(order) -> []Trade
- GetOrderBook(shareID) -> OrderBook

ShareService:
- GetShare(shareID) -> Share
- GetSharePrice(shareID) -> float64
- UpdateSharePrice(shareID, price)
- GetOrderBook(shareID) -> OrderBook

PortfolioService:
- GetPortfolio(userID) -> []Holding
- GetHoldings(userID, shareID) -> Holding
- UpdateHoldings(userID, shareID, quantity, price)
- ValidateSell(userID, shareID, quantity) -> bool

OrderPendingService:
- AddOrder(order)
- RemoveOrder(orderID)
- GetBestBuyOrder(shareID) -> *Order
- GetBestSellOrder(shareID) -> *Order
- GetOrdersByPrice(shareID, operation, price) -> []Order

ADDITIONAL CONSIDERATIONS:
1. Matching Priority:
   - Price Priority: Best price first
   - Time Priority: First come first serve (same price)

2. Order Types:
   - LIMIT: Execute at specific price or better
   - MARKET: Execute at best available price immediately

3. Error Handling:
   - Insufficient balance for buy
   - Insufficient shares for sell
   - Invalid price/quantity
   - Order not found/already executed

4. Concurrency:
   - Thread-safe order matching
   - Lock on order book per share
   - Atomic trade execution

5. Data Structures:
   - Use priority queues (heap) for O(log N)
   - Buy orders: Max heap (highest price first)
   - Sell orders: Min heap (lowest price first)

6. Additional Features:
   - Order history
   - Trade history
   - Market depth (order book levels)
   - Price history

ADDITIONAL DESIGN PATTERNS:
- State Pattern: Order status transitions
- Repository Pattern: Data access abstraction

*/