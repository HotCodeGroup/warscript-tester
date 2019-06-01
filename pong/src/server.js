const port = 5000

const loadcode = "/loadcode"
const run = "/run"

let ready = false
let p = new Function("me", "enemy", "ball", "game", "console", "{}");

const http = require('http');
const {
    parse
} = require('querystring');

const server = http.createServer((req, res) => {
    console.log(req.url)
    if (req.method === 'POST' && req.url === loadcode) {
        let body = [];
        req.on('data', (chunk) => {
            body.push(chunk);
        }).on('end', () => {
            body = Buffer.concat(body).toString();
            let params = JSON.parse(body)

            let errorText = ""
            try {
                p = new Function("me", "enemy", "ball", "game", "console", params.code);
            } catch (err) {
                errorText = `ERROR: ${err.name}; ${err.message};`
            }
            let resp = JSON.stringify({
                "error": errorText,
            })
            res.setHeader('Content-Type', 'application/json');
            res.statusCode = 200;
            res.end(resp)
        });
    } else if (req.method === 'POST' && req.url === run) {
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
            me.height = params.me.height
            me.width = params.me.width
            me.setMoveVector = function(speed, x, y) {
                let nSpeed = speed / Math.sqrt(x * x + y * y);
                if (isNaN(nSpeed) || nSpeed == Infinity) {
                    nSpeed = 0;
                }

                this.vX = x * nSpeed;
                this.vY = y * nSpeed;
            }

            let enemy = {}
            enemy.x = params.enemy.x
            enemy.y = params.enemy.y
            enemy.vX = params.enemy.vX
            enemy.vY = params.enemy.vY
            enemy.height = params.enemy.height
            enemy.width = params.enemy.width

            let ball = {}
            ball.x = params.ball.x
            ball.y = params.ball.y
            ball.vX = params.ball.vX
            ball.vY = params.ball.vY
            ball.height = params.ball.height
            ball.width = params.ball.width

            let game = {}
            game.height = params.game.height
            game.width = params.game.width
            game.ticks_left = params.game.ticks_left

            let cnsl = {}
            cnsl.logs = []
            cnsl.log = function(msg) {
                this.logs.push(msg)
            }


            let errorText = ""
            try {
                p(me, enemy, ball, game, cnsl)
            } catch (err) {
                errorText = `[TICK: ${10000 - params.game.ticks_left}] ERROR: ${err.name}; ${err.message};`
                cnsl.logs.push(`[TICK: ${10000 - params.game.ticks_left}] ERROR: ${err.name}; ${err.message};`)
            }

            let resp = JSON.stringify({
                "me": me,
                "enemy": enemy,
                "ball": ball,
                "console": cnsl.logs,
                "error": errorText
            })
            res.setHeader('Content-Type', 'application/json');
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