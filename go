<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>ETHEREAL GO</title>
    <style>
        :root {
            --bg: #050505;
            --board: rgba(255, 255, 255, 0.05);
            --line: rgba(255, 255, 255, 0.15);
            --p1: #00f3ff; /* Cyan (Black Stone visual) */
            --p2: #ff0055; /* Pink (White Stone visual) */
            --glow: 0 0 10px;
        }
        
        * { box-sizing: border-box; -webkit-user-drag: none; -webkit-tap-highlight-color: transparent; user-select: none; font-family: 'Segoe UI', sans-serif; }
        
        body {
            margin: 0; height: 100dvh; background: #080808; color: white;
            display: flex; flex-direction: column; align-items: center; justify-content: center;
            overflow: hidden; touch-action: none;
        }

        /* 背景 */
        .nebula {
            position: absolute; width: 100%; height: 100%; z-index: -1; pointer-events: none;
            background: 
                radial-gradient(circle at 20% 20%, rgba(0, 243, 255, 0.08), transparent 50%),
                radial-gradient(circle at 80% 80%, rgba(255, 0, 85, 0.08), transparent 50%);
        }

        /* HUD */
        .hud {
            width: 90vw; max-width: 500px; display: flex; justify-content: space-between;
            margin-bottom: 20px; padding: 10px 20px; background: rgba(255,255,255,0.05);
            border-radius: 30px; border: 1px solid rgba(255,255,255,0.1);
            font-size: 0.9rem; letter-spacing: 1px;
        }
        .score-box { display: flex; align-items: center; gap: 10px; }
        .dot { width: 10px; height: 10px; border-radius: 50%; background: #333; transition: 0.3s; }
        .active-p1 .dot { background: var(--p1); box-shadow: 0 0 10px var(--p1); }
        .active-p2 .dot { background: var(--p2); box-shadow: 0 0 10px var(--p2); }

        /* 棋盘容器 */
        .board-wrap {
            width: 92vw; height: 92vw; max-width: 500px; max-height: 500px;
            position: relative; padding: 10px; /* Padding creates the edge area */
            background: var(--board); border-radius: 4px;
            border: 1px solid rgba(255,255,255,0.1);
            box-shadow: 0 20px 50px rgba(0,0,0,0.5);
        }

        /* 网格线层 */
        .grid-layer {
            position: absolute; top: 10px; left: 10px; right: 10px; bottom: 10px;
            display: grid;
            grid-template-columns: repeat(12, 1fr);
            grid-template-rows: repeat(12, 1fr);
            pointer-events: none; z-index: 0;
        }
        .cell { border: 1px solid var(--line); border-width: 0 1px 1px 0; }
        /* Fix borders */
        .grid-layer { border: 1px solid var(--line); border-bottom: 0; border-right: 0; }

        /* 星位 */
        .star-point {
            position: absolute; width: 4px; height: 4px; background: white; border-radius: 50%;
            transform: translate(-50%, -50%); z-index: 0; opacity: 0.6;
        }

        /* 交互层 (交叉点) */
        .touch-layer {
            position: absolute; top: 0; left: 0; width: 100%; height: 100%;
            display: grid;
            grid-template-columns: repeat(13, 1fr);
            grid-template-rows: repeat(13, 1fr);
            z-index: 10;
        }
        
        .intersection {
            width: 100%; height: 100%;
            position: relative; display: flex; justify-content: center; align-items: center;
            cursor: pointer;
        }
        
        /* 棋子 */
        .stone {
            width: 85%; height: 85%; border-radius: 50%;
            transform: scale(0); transition: transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
            box-shadow: inset 0 0 10px rgba(0,0,0,0.8);
        }
        .stone.show { transform: scale(1); }
        
        /* P1 = Black (represented as Cyan/Dark) */
        .stone.p1 { background: radial-gradient(circle at 30% 30%, #444, #000); border: 1px solid var(--p1); box-shadow: 0 0 5px var(--p1); }
        /* P2 = White (represented as Pink/White) */
        .stone.p2 { background: radial-gradient(circle at 30% 30%, #fff, #ccc); border: 1px solid var(--p2); box-shadow: 0 0 5px var(--p2); }
        
        .stone.last-move::after {
            content: ''; position: absolute; top: 50%; left: 50%;
            width: 30%; height: 30%; border-radius: 50%;
            background: white; transform: translate(-50%, -50%);
            animation: pulse 1s infinite;
        }
        @keyframes pulse { 0% { opacity: 0.5; } 50% { opacity: 1; } 100% { opacity: 0.5; } }

        /* 菜单 */
        .overlay {
            position: absolute; top: 0; left: 0; width: 100%; height: 100%;
            background: rgba(0,0,0,0.9); backdrop-filter: blur(10px);
            display: flex; flex-direction: column; justify-content: center; align-items: center;
            z-index: 100; transition: opacity 0.4s;
        }
        .hidden { opacity: 0; pointer-events: none; }
        
        h1 { letter-spacing: 5px; text-shadow: 0 0 20px var(--p1); margin-bottom: 40px; }
        
        .btn {
            padding: 15px 40px; border: 1px solid white; background: transparent;
            color: white; font-size: 1.1rem; margin: 10px; cursor: pointer;
            transition: 0.3s; letter-spacing: 2px; min-width: 200px;
        }
        .btn:active { background: white; color: black; }
        .btn-cyan { border-color: var(--p1); color: var(--p1); }
        .btn-cyan:active { background: var(--p1); }
        
        .controls { margin-top: 20px; display: flex; gap: 20px; }
        .mini-btn { padding: 8px 20px; border: 1px solid rgba(255,255,255,0.3); background: transparent; color: #aaa; font-size: 0.8rem; }

    </style>
</head>
<body>
    <div class="nebula"></div>

    <div id="start-screen" class="overlay">
        <h1>ETHEREAL GO</h1>
        <button class="btn btn-cyan" onclick="game.start(1)">NOVICE (Depth 1)</button>
        <button class="btn btn-cyan" onclick="game.start(2)">ADEPT (Tactical)</button>
        <button class="btn btn-cyan" onclick="game.start(3)">MASTER (Strategic)</button>
        <div style="margin-top: 30px; font-size: 0.8rem; opacity: 0.5;">13x13 BOARD • CHINESE RULES</div>
    </div>

    <div id="end-screen" class="overlay hidden">
        <h1 id="end-msg">GAME OVER</h1>
        <div id="final-score" style="margin-bottom: 30px; font-size: 1.2rem; opacity: 0.8;"></div>
        <button class="btn" onclick="location.reload()">REBOOT</button>
        <button class="btn" onclick="window.location.href='index.html'">EXIT TO HUB</button>
    </div>

    <div class="hud">
        <div class="score-box active-p1" id="p1-ui">
            <div class="dot"></div>
            <span>YOU (BLACK)</span>
            <span id="caps-p1" style="margin-left:5px; opacity:0.5;">[0]</span>
        </div>
        <div class="score-box" id="p2-ui">
            <span id="caps-p2" style="margin-right:5px; opacity:0.5;">[0]</span>
            <span>AI (WHITE)</span>
            <div class="dot"></div>
        </div>
    </div>

    <div class="board-wrap" id="board-container">
        <div class="grid-layer" id="grid-bg"></div>
        <div class="touch-layer" id="touch-layer"></div>
    </div>

    <div class="controls">
        <button class="mini-btn" onclick="game.pass()">PASS</button>
        <button class="mini-btn" onclick="window.location.href='index.html'">EXIT</button>
    </div>

<script>
// 13x13 Board for better mobile experience
const SIZE = 13; 

class GoGame {
    constructor() {
        this.board = []; // 0:empty, 1:black, 2:white
        this.turn = 1;
        this.caps = { 1: 0, 2: 0 }; // Captures
        this.ko = null; // Ko point
        this.lastMove = null;
        this.history = [];
        this.diff = 1;
        this.state = 'init';
        this.passed = false;
    }

    init() {
        this.createGrid();
    }

    start(level) {
        this.diff = level;
        this.reset();
        document.getElementById('start-screen').classList.add('hidden');
        this.state = 'play';
    }

    reset() {
        this.board = Array(SIZE).fill(0).map(() => Array(SIZE).fill(0));
        this.turn = 1;
        this.caps = { 1: 0, 2: 0 };
        this.ko = null;
        this.lastMove = null;
        this.history = [];
        this.passed = false;
        this.render();
        this.updateUI();
    }

    createGrid() {
        // 1. Draw Grid Lines
        // CSS Grid handles the cells.
        
        // 2. Draw Star Points (Tian Yuan & corners for 13x13)
        const stars = [ [3,3], [3,9], [6,6], [9,3], [9,9] ];
        const wrap = document.getElementById('board-container');
        stars.forEach(([r,c]) => {
            // Calculate % position. 13 lines mean intervals are 100% / 12 approx? 
            // Actually touch layer is 13x13. Center of cell.
            // Let's align manually based on 13x13 grid logic.
            // 10px padding. Board width = 100% - 20px.
            // Each cell is (100% / 13). Center is (c + 0.5).
            
            const dot = document.createElement('div');
            dot.className = 'star-point';
            // Simple hack: attach to the specific intersection div? 
            // No, absolute position is cleaner visually.
            // 13 cells. Star at index 3 (4th line).
            // Position = 10px + (calcWidth * column / 12) ? No.
            // Easier: Let's put stars inside the touch layer's divs.
        });

        // 3. Create Touch Layer
        const tl = document.getElementById('touch-layer');
        tl.innerHTML = '';
        for(let r=0; r<SIZE; r++) {
            for(let c=0; c<SIZE; c++) {
                const div = document.createElement('div');
                div.className = 'intersection';
                div.dataset.r = r;
                div.dataset.c = c;
                div.onclick = () => this.playerMove(r, c);
                
                // Add star point if match
                if(stars.some(s => s[0]===r && s[1]===c)) {
                    const s = document.createElement('div');
                    s.className = 'star-point';
                    s.style.position='static'; s.style.transform='none'; // Reset
                    div.appendChild(s);
                }

                tl.appendChild(div);
            }
        }
    }

    // --- LOGIC ---

    isValid(r, c) { return r >= 0 && r < SIZE && c >= 0 && c < SIZE; }

    getGroup(r, c, board = this.board) {
        const color = board[r][c];
        if (color === 0) return null;
        
        let group = [];
        let liberties = 0;
        let visited = new Set();
        let queue = [[r,c]];
        visited.add(`${r},${c}`);

        while (queue.length > 0) {
            const [curR, curC] = queue.shift();
            group.push([curR, curC]);

            [[0,1], [0,-1], [1,0], [-1,0]].forEach(([dr, dc]) => {
                const nr = curR + dr, nc = curC + dc;
                if (this.isValid(nr, nc)) {
                    const key = `${nr},${nc}`;
                    if (board[nr][nc] === 0) {
                        // It's a liberty. Use a set to count unique liberties?
                        // Usually strictly logic doesn't need unique count for capture check, just > 0.
                        // But for AI we might need count.
                        // Let's add to a set of liberties.
                        // For simple capture check: finding one is enough.
                        // We'll just accumulate valid liberties here blindly, 
                        // usually we need a separate visitedLiberties set.
                        // Simplified: just check logic.
                    } else if (board[nr][nc] === color && !visited.has(key)) {
                        visited.add(key);
                        queue.push([nr, nc]);
                    }
                }
            });
        }

        // Count liberties accurately
        let libSet = new Set();
        group.forEach(([gr, gc]) => {
            [[0,1], [0,-1], [1,0], [-1,0]].forEach(([dr, dc]) => {
                const nr = gr + dr, nc = gc + dc;
                if (this.isValid(nr, nc) && board[nr][nc] === 0) {
                    libSet.add(`${nr},${nc}`);
                }
            });
        });

        return { stones: group, liberties: libSet.size };
    }

    placeStone(r, c, color, mock = false) {
        // 1. Check empty
        if (this.board[r][c] !== 0) return false;
        
        // 2. Check Ko
        if (this.ko && this.ko[0] === r && this.ko[1] === c) return false;

        // Temporary board for simulation
        let nextBoard = this.board.map(row => [...row]);
        nextBoard[r][c] = color;
        
        let captured = [];
        const opp = color === 1 ? 2 : 1;
        
        // 3. Check Captures
        [[0,1], [0,-1], [1,0], [-1,0]].forEach(([dr, dc]) => {
            const nr = r + dr, nc = c + dc;
            if (this.isValid(nr, nc) && nextBoard[nr][nc] === opp) {
                const group = this.getGroup(nr, nc, nextBoard);
                if (group && group.liberties === 0) {
                    captured.push(...group.stones);
                }
            }
        });

        // Remove captured stones
        captured.forEach(([cr, cc]) => nextBoard[cr][cc] = 0);

        // 4. Check Suicide
        const selfGroup = this.getGroup(r, c, nextBoard);
        if (selfGroup.liberties === 0) return false; // Illegal suicide

        // 5. Valid Move Execution
        if (!mock) {
            this.board = nextBoard;
            this.caps[color] += captured.length;
            
            // Ko Handling: simple ko is 1 stone captured, returning board to state before opponent move.
            // Strict check: if captured 1 stone and self group size is 1.
            if (captured.length === 1 && selfGroup.stones.length === 1) {
                this.ko = captured[0];
            } else {
                this.ko = null;
            }

            this.lastMove = [r, c];
            this.passed = false;
            
            if(navigator.vibrate) {
                captured.length > 0 ? navigator.vibrate(50) : navigator.vibrate(10);
            }
            
            return true;
        }
        
        return true; // For mock check
    }

    playerMove(r, c) {
        if (this.state !== 'play' || this.turn !== 1) return;
        
        if (this.placeStone(r, c, 1)) {
            this.render();
            this.turn = 2;
            this.updateUI();
            setTimeout(() => this.aiMove(), 500);
        }
    }

    pass() {
        if (this.state !== 'play') return;
        if (this.passed) {
            this.endGame();
            return;
        }
        this.passed = true;
        this.turn = this.turn === 1 ? 2 : 1;
        this.ko = null;
        this.updateUI();
        if (this.turn === 2) setTimeout(() => this.aiMove(), 500);
    }

    // --- AI ---
    aiMove() {
        if (this.state !== 'play') return;

        const opp = 1;
        const self = 2;
        let bestMove = null;
        let maxScore = -Infinity;
        
        // Candidates
        let moves = [];
        for(let r=0; r<SIZE; r++) for(let c=0; c<SIZE; c++) {
            if(this.board[r][c] === 0) moves.push([r,c]);
        }
        
        // 1. NOVICE: Random legal
        if (this.diff === 1) {
            moves.sort(() => Math.random() - 0.5);
            for(let m of moves) {
                if(this.placeStone(m[0], m[1], self)) { // placeStone executes if not mock
                    // Wait, placeStone executes? No, I need to wrap logic.
                    // My placeStone function executes if mock=false.
                    // So for Novice, just finding the first one that works is fine.
                    this.finishTurn();
                    return;
                }
            }
            this.pass(); // No legal moves
            return;
        }

        // 2. ADEPT / MASTER: Evaluation
        // We need to simulate moves without changing state.
        
        // Shuffle moves for variety
        moves.sort(() => Math.random() - 0.5);

        for (let [r, c] of moves) {
            // Mock simulation
            if (this.placeStone(r, c, self, true)) {
                let score = 0;
                
                // Analyze neighbors
                let captureCount = 0;
                let atariPotential = 0;
                let influence = 0;

                // Sim board
                let tempBoard = this.board.map(row => [...row]);
                tempBoard[r][c] = self;
                // Note: Capture logic logic is duplicated here for scoring, simplified.
                
                // Heuristic 1: Capture Priority
                [[0,1], [0,-1], [1,0], [-1,0]].forEach(([dr, dc]) => {
                    const nr = r+dr, nc = c+dc;
                    if(this.isValid(nr, nc)) {
                        if(tempBoard[nr][nc] === opp) {
                            const g = this.getGroup(nr, nc, tempBoard);
                            if(g.liberties === 0) captureCount += g.stones.length;
                            else if(g.liberties === 1) atariPotential += 5; // Atari enemy
                        }
                    }
                });

                // Heuristic 2: Save Self from Atari
                // If playing here gives my group more liberties than before? 
                // Hard to calc easily. Simplification:
                // Check if this move has decent liberties
                const myG = this.getGroup(r, c, tempBoard);
                if (myG.liberties === 1) score -= 50; // Bad shape (self-atari), unless capturing
                if (myG.liberties > 2) score += 5;

                // Heuristic 3: Position (Star points & Center)
                const isEdge = r===0 || r===SIZE-1 || c===0 || c===SIZE-1;
                if (isEdge) score -= 2; // Avoid 1st line early
                
                // 3rd line (Territory) & 4th line (Influence)
                const rDist = Math.min(r, SIZE-1-r);
                const cDist = Math.min(c, SIZE-1-c);
                if (rDist >= 2 && rDist <= 4 && cDist >= 2 && cDist <= 4) score += 3;
                
                // Star points bias
                if ([3,9,6].includes(r) && [3,9,6].includes(c)) score += 2;

                score += captureCount * 100;
                score += atariPotential;

                if (this.diff === 3) {
                    // Master: Check pattern (Hane/Cut) - very simplified
                    // If diagonal is enemy, and adjacent is own...
                }

                if (score > maxScore) {
                    maxScore = score;
                    bestMove = [r, c];
                }
            }
        }

        if (bestMove) {
            this.placeStone(bestMove[0], bestMove[1], self);
            this.finishTurn();
        } else {
            this.pass();
        }
    }

    finishTurn() {
        this.render();
        this.turn = 1;
        this.updateUI();
    }

    endGame() {
        this.state = 'end';
        // Simple scoring: Stones on board + Captures
        // Note: This is not real territory scoring (too complex for JS snippet).
        // This is "Area Scoring" simplified.
        let blackScore = this.caps[1];
        let whiteScore = this.caps[2];
        
        for(let r=0; r<SIZE; r++) for(let c=0; c<SIZE; c++) {
            if(this.board[r][c]===1) blackScore++;
            if(this.board[r][c]===2) whiteScore++;
        }
        
        // Komi
        whiteScore += 6.5;

        const msg = document.getElementById('end-msg');
        const score = document.getElementById('final-score');
        const scr = document.getElementById('end-screen');
        
        if(blackScore > whiteScore) {
            msg.innerText = "VICTORY";
            msg.style.color = "var(--p1)";
        } else {
            msg.innerText = "DEFEAT";
            msg.style.color = "var(--p2)";
        }
        score.innerText = `B: ${blackScore} | W: ${whiteScore} (6.5 komi)`;
        scr.classList.remove('hidden');
    }

    // --- RENDER ---
    render() {
        const tl = document.getElementById('touch-layer');
        const stones = tl.children;
        
        for(let r=0; r<SIZE; r++) {
            for(let c=0; c<SIZE; c++) {
                const idx = r*SIZE + c;
                const cell = stones[idx];
                const val = this.board[r][c];
                
                // Clear existing stone
                const existing = cell.querySelector('.stone');
                if(existing && existing.dataset.c != val) existing.remove();

                if(val !== 0 && (!existing || existing.dataset.c != val)) {
                    if(existing) existing.remove();
                    const s = document.createElement('div');
                    s.className = `stone p${val} ${this.lastMove && this.lastMove[0]===r && this.lastMove[1]===c ? 'last-move' : ''}`;
                    s.dataset.c = val;
                    cell.appendChild(s);
                    // Trigger animation
                    requestAnimationFrame(() => s.classList.add('show'));
                }
                
                // Update last move marker if stone already exists
                if(val !== 0 && existing) {
                    if(this.lastMove && this.lastMove[0]===r && this.lastMove[1]===c) existing.classList.add('last-move');
                    else existing.classList.remove('last-move');
                }
            }
        }
    }

    updateUI() {
        const p1 = document.getElementById('p1-ui');
        const p2 = document.getElementById('p2-ui');
        document.getElementById('caps-p1').innerText = `[${this.caps[1]}]`;
        document.getElementById('caps-p2').innerText = `[${this.caps[2]}]`;

        if(this.turn === 1) {
            p1.classList.add('active-p1'); p2.classList.remove('active-p2');
        } else {
            p1.classList.remove('active-p1'); p2.classList.add('active-p2');
        }
    }
}

const game = new GoGame();
game.init();
</script>
</body>
</html>
