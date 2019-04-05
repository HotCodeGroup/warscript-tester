const port = 5000

const loadcode = "/loadcode"
const run = "/run"

let ready = false
let p = new Function("me", "enemy", "ball", "{}");

const http = require('http');
const {
    parse
} = require('querystring');

const server = http.createServer((req, res) => {
    console.log(req.url)
    if (req.method === 'POST' && req.url == loadcode) {
        let body = [];
        req.on('data', (chunk) => {
            body.push(chunk);
        }).on('end', () => {
            body = Buffer.concat(body).toString();
            let params = JSON.parse(body)
            p = new Function("me", "enemy", "ball", params.code);
            res.setHeader("Content-Type", "application/json");
            res.status= 200;
            res.end(params.code)
        });
    } else if (req.method === 'POST' && req.url == run) {
        let body = [];
        req.on('data', (chunk) => {
            body.push(chunk);
            console.log(chunk);
        }).on('end', () => {
            body = Buffer.concat(body).toString();
            let params = JSON.parse(body)
            let me = {}
            me.x = params.me.x
            me.y = params.me.y
            me.vX = params.me.vX
            me.vY = params.me.vY

            let enemy = {}
            enemy.x = params.enemy.x
            enemy.y = params.enemy.y
            enemy.vX = params.enemy.vX
            enemy.vY = params.enemy.vY

            let ball = {}
            ball.x = params.ball.x
            ball.y = params.ball.y
            ball.vX = params.ball.vX
            ball.vY = params.ball.vY

            me.setMoveVector = (speed, x, y) => {

                let nSpeed = speed / Math.sqrt(x * x + y * y);
                if (isNaN(nSpeed) || nSpeed == Infinity) {
                    nSpeed = 0;
                }

                this.vX = x * nSpeed;
                this.vY = y * nSpeed;
            }
            p(me, enemy, ball)

            let resp = JSON.stringify({"me":me, "enemy":enemy, "ball":ball})
            res.setHeader("Content-Type", "application/json");
            res.status(200);
            res.end(resp)
        });
    }
});
server.listen(port);

function collectRequestData(request, callback) {
    const FORM_URLENCODED = 'application/x-www-form-urlencoded';
    if (request.headers['content-type'] === FORM_URLENCODED) {
        let body = '';
        request.on('data', chunk => {
            body += chunk.toString();
        });
        request.on('end', () => {
            callback(parse(body));
        });
    } else {
        callback(null);
    }
}