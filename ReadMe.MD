**General Ledger(GL)**
----
Introduction
============
Design and implement general ledger with deposit, withdrew and query of transactions.

Http Rest API, Html and java script test page

Golang: ver 1.18

PostgresSQL: ver 14.1

Main features
- JWT support for security with configurable duration. 
- Service time out support with context and channel
- Dictionary result set
- Multi media log support
- Database connection pool with refresh period time
- Encrypt or plane start configuration 


StartApplication
============
- Configuration file

You can find config file in "./Files/ConfigPlane.json" and change start IP, Port and change postgresSQL database config on your own DB.

- Run
    - go run ./main.go
- Build
    - go build ./main.go


Services
============

WebPage
============
- Web page application
    - /page: Login, create transaction and view transactions

Signing
============
  
- {IP}:{PORT}/login: Takes Username and Password with json format. [POST] 
  - Request
```json
  {
  "username":"mohsen",
  "password":"qaz@123456"
  }
```
- Response
  - Return status, time, valid duration, account numbers and sing key for reuse in other services

```json
{
  "BStatus": 1,
  "Time": "2022-06-05T13:01:27.7407569+01:00",
  "Result": {
    "Username": "mohsen",
    "IsActive": true,
    "ValidTill": "Sun Jun 5 13:04:27 BST 2022",
    "Sign": "nzTTN4nD7U4kdar7KkHOx9wOPd8zTcoV7eDCZaiITTE/mD//E02QVIcWVno2oWIToCwZPxy/Vot/yqfNmwXr",
    "Accounts": [
      {
        "id": "76ccfbc5-438d-4310-a807-01424d39202d",
        "userId": "a061d6f0-e119-11ec-8c31-f8a9634eef20",
        "accNo": "914300118110",
        "balance": 101.22
      }
    ]
  }
}
```
Create transaction(deposit, withdrew)
============
- {IP}:{PORT}/transaction [POST]
  - Takes account number, amount, note, transaction type (deposit or withdrew) and sign key.
```json
{
  "transaction": {
    "account": "914300118110",
    "amount": 1,
    "note": "test",
    "action": 1
  },
  "sign": "PoyxpMRJLcP/+81puXsRAK0zAgEYzYpovWY/ZZzGvgTj5VFIhwaR5fEewV2Wenk7Hs67BTK3PXLkIMvLiN+6"
}
```
- Response
    - Returns status and result description
```json
{
  "BStatus": 1,
  "Time": "2022-06-05T20:22:03.976055+01:00",
  "Result": "transaction commit"
}
```
Transaction Query
============
- {IP}:{PORT}/transactions [POST]
  - Takes account number, start date, end date and sign key.

- Request
```json
{
  "account": "914300118110",
  "fromDate": "2022-06-05",
  "toDate": "2022-06-05",
  "sign": "Hq9mSTP95EcExMeG5siicHedLB2Yaq+Cjgl27mUViu4dyFExm/e6zx8rtiMKzG9BjA6nd5MDODE1PI0Hnf3F"
}
```

- Response
```json
{
  "BStatus": 1,
  "Time": "2022-06-05T20:31:47.4779118+01:00",
  "Result": [
    {
      "account": "914300118110",
      "amount": 2,
      "note": "test-1",
      "action": 0
    },
    {
      "account": "914300118110",
      "amount": -12345678,
      "note": "test-2",
      "action": 0
    },
    {
      "account": "914300118110",
      "amount": 100,
      "note": "test-3",
      "action": 0
    },
    {
      "account": "914300118110",
      "amount": 1.22,
      "note": "test-4",
      "action": 0
    },
    {
      "account": "914300118110",
      "amount": 1,
      "note": "test-5",
      "action": 0
    }
  ]
}
```