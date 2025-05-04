# xyz-kredit-plus
This project is created regarding to <strong>Kredit Plus Technical Test</strong>. \
Problem definitions and requirements: [PRD.pdf](./PRD.pdf)

## How to Run
1. Clone repository
```
git clone https://github.com/michaelyusak/xyz-kredit-plus.git
cd xyz-kredit-plus
```
2. Install dependencies
```
go mod tidy
```
3. Create `config.json` and `.mysql.env`. See [config.example.json](./config.example.json) and [mysql.example.env](./mysql.example.env).
4. Run App
```
docker-compose-up
```
5. Check documentation. Clink [here](http://localhost:8080). Or go to
```
{HOST}/swagger/index.html
```
## Adjustments
Due to the lack of technical information, here are several adjustment applied on this app.
### Flow
```
Register -> Process KYC -> Login -> Create Transaction
```
First, user has to own an account, tokens will be granted if register succeeded. Then, no need to login, user has to undergo KYC process to submit consumer data. After KYC completed, since there is no session engine yet, user has to login to create a new access token. Finaly, user can create transaction.

### Remaining Limit Calculation
Since there is no given information about how limit is calculated, here is how limit is calculated in this app.
#### Find the Discount
If user requested a transaction with 2 month installemnt and Rp100000 OTR, the discount is as shown below.
```
disc = (limit_2_m - OTR) / limit_2_m
```
#### Adjust Limits
The limit of 1, 2, 3, and 4 months installemnt will be recalculated by multiplying the discount to the current limit.
```
newLimit = currLimit * disc
```