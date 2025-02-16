const ff = @import("ff");

const draw = ff.draw;

var buf: [1735]u8 = undefined;
var fff: ff.Font = undefined;
var btn: ff.Buttons = undefined;
var pad: ff.Pad = undefined;
var pal: ff.Palette = .{
    .black = 0x000000,
    .gray = 0x292929,
    .white = 0xffffff,
    .orange = 0xf7a41d,
};

pub export fn boot() void {
    pal.set();

    fff = ff.loadFile("font", buf[0..]);
}

pub export fn update() void {
    const me = ff.getMe();

    btn = ff.readButtons(me);
    pad = ff.readPad(me).?;
}

pub export fn render() void {
    ff.clearScreen(.black);

    renderPad();
    renderButtons(.{ .x = 155, .y = 100 });
    renderZigLogo(.{ .x = 3 + @divExact(pad.x, 500), .y = 6 + -@divExact(pad.y, 500) });
}

fn renderZigLogo(offset: ff.Point) void {
    const x = offset.x;
    const y = offset.y;
    const w = ff.Style{ .fill_color = .white };
    const o = ff.Style{ .fill_color = .orange };
    const b = ff.Style{ .fill_color = .black };

    { // [⚡]
        draw.rect(x + 2, y + 13, 92, 57, o);
        draw.tri(x + 20, y + 27, x + 30, y + 13, x + 24, y + 28, b);
        draw.tri(x + 30, y + 13, x + 24, y + 28, x + 35, y + 13, b);
        draw.tri(x + 21, y + 57, x + 10, y + 70, x + 15, y + 70, b);
        draw.tri(x + 21, y + 57, x + 15, y + 70, x + 26, y + 57, b);
        draw.rect(x + 17, y + 27, 66, 31, b);
        draw.tri(x + 5, y + 83, x + 51, y + 26, x + 45, y + 56, o);
        draw.tri(x + 45, y + 56, x + 53, y + 16, x + 89, y + 0, o);
    }
    { // Z
        draw.tri(x + 143, y + 22, x + 109, y + 61, x + 126, y + 60, w);
        draw.tri(x + 126, y + 60, x + 141, y + 22, x + 159, y + 24, w);
        draw.rect(x + 109, y + 13, 50, 11, w);
        draw.rect(x + 109, y + 60, 49, 11, w);
    }
    { // I
        draw.rect(x + 166, y + 13, 13, 58, w);
    }
    { // G
        draw.tri(x + 75, y + 27, x + 86, y + 13, x + 79, y + 13, b);
        draw.tri(x + 79, y + 13, x + 66, y + 31, x + 75, y + 27, b);
        draw.tri(x + 78, y + 57, x + 66, y + 72, x + 72, y + 57, b);
        draw.tri(x + 72, y + 57, x + 67, y + 70, x + 62, y + 70, b);
        draw.elli(x + 184, y + 14, 56, 60, w);
        draw.elli(x + 197, y + 27, 32, 34, b);
        draw.tri(x + 209, y + 41, x + 237, y + 53, x + 237, y + 41, w);
        draw.tri(x + 241, y + 17, x + 219, y + 40, x + 241, y + 40, b);
        draw.rect(x + 223, y + 41, 10, 23, w);
        draw.rect(x + 233, y + 41, 10, 23, b);
    }
}

fn renderButtons(at: ff.Point) void {
    draw.Text("WEST", fff, .{ .x = at.x + 0, .y = at.y + 0 }, if (btn.w) .orange else .gray);
    draw.Text("SOUTH", fff, .{ .x = at.x + 0, .y = at.y + 20 }, if (btn.s) .orange else .gray);
    draw.Text("NORTH", fff, .{ .x = at.x + 40, .y = at.y + 0 }, if (btn.n) .orange else .gray);
    draw.Text("EAST", fff, .{ .x = at.x + 40, .y = at.y + 20 }, if (btn.e) .orange else .gray);
}

fn renderPad() void {
    const g = ff.Style{ .fill_color = .gray };
    const o = ff.Style{ .fill_color = .orange };
    const px = 50 + @divExact(pad.x, 50);
    const py = 100 - @divExact(pad.y, 50);

    draw.line(0, 0, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    draw.line(240, 0, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    draw.line(0, 160, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    draw.line(240, 160, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    draw.circ(px, py, 16, if (pad.x + pad.y == 0) g else o);
}
