const ff = @import("ff");

var fff: ff.Font = undefined;
var btn: ff.Buttons = undefined;
var pad: ff.Pad = undefined;

pub export fn boot() void {
    ff.setColorHex(.black, 0x000000);
    ff.setColorHex(.gray, 0x292929);
    ff.setColorHex(.white, 0xffffff);
    ff.setColorHex(.orange, 0xf7a41d);

    var buf: [1735]u8 = undefined;

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
    renderZigLogo(.{
        .x = 3 + @divExact(pad.x, 500),
        .y = 6 + -@divExact(pad.y, 500),
    });
    renderButtons();
}

fn renderZigLogo(offset: ff.Point) void {
    const x = offset.x;
    const y = offset.y;
    const w = ff.Style{ .fill_color = .white };
    const o = ff.Style{ .fill_color = .orange };
    const b = ff.Style{ .fill_color = .black };

    rect(x + 2, y + 13, 92, 57, o);
    rect(x + 109, y + 13, 50, 11, w);
    rect(x + 109, y + 60, 49, 11, w);
    rect(x + 166, y + 13, 13, 58, w);
    tri(x + 143, y + 22, x + 109, y + 61, x + 126, y + 60, w);
    tri(x + 126, y + 60, x + 141, y + 22, x + 159, y + 24, w);
    tri(x + 20, y + 27, x + 30, y + 13, x + 24, y + 28, b);
    tri(x + 30, y + 13, x + 24, y + 28, x + 35, y + 13, b);
    tri(x + 21, y + 57, x + 10, y + 70, x + 15, y + 70, b);
    tri(x + 21, y + 57, x + 15, y + 70, x + 26, y + 57, b);
    rect(x + 17, y + 27, 66, 31, b);
    tri(x + 5, y + 83, x + 51, y + 26, x + 45, y + 56, o);
    tri(x + 45, y + 56, x + 53, y + 16, x + 89, y + 0, o);
    tri(x + 75, y + 27, x + 86, y + 13, x + 79, y + 13, b);
    tri(x + 79, y + 13, x + 66, y + 31, x + 75, y + 27, b);
    tri(x + 78, y + 57, x + 66, y + 72, x + 72, y + 57, b);
    tri(x + 72, y + 57, x + 67, y + 70, x + 62, y + 70, b);
    elli(x + 184, y + 14, 56, 60, w);
    elli(x + 197, y + 27, 32, 34, b);
    tri(x + 209, y + 41, x + 237, y + 53, x + 237, y + 41, w);
    tri(x + 241, y + 17, x + 219, y + 40, x + 241, y + 40, b);
    rect(x + 223, y + 41, 10, 23, w);
    rect(x + 233, y + 41, 10, 23, b);
}

fn renderButtons() void {
    ff.drawText("WEST", fff, .{ .x = 160, .y = 100 }, if (btn.w) .orange else .gray);
    ff.drawText("SOUTH", fff, .{ .x = 160, .y = 120 }, if (btn.s) .orange else .gray);
    ff.drawText("NORTH", fff, .{ .x = 200, .y = 100 }, if (btn.n) .orange else .gray);
    ff.drawText("EAST", fff, .{ .x = 200, .y = 120 }, if (btn.e) .orange else .gray);
}

fn renderPad() void {
    const g = ff.Style{ .fill_color = .gray };
    const o = ff.Style{ .fill_color = .orange };
    const px = 50 + @divExact(pad.x, 50);
    const py = 100 - @divExact(pad.y, 50);

    line(0, 0, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    line(240, 0, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    line(0, 160, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    line(240, 160, px + 8, py + 8, .{ .color = .gray, .width = 2 });
    circ(px, py, 16, if (pad.x + pad.y == 0) g else o);
}

fn tri(x1: i32, y1: i32, x2: i32, y2: i32, x3: i32, y3: i32, s: ff.Style) void {
    ff.drawTriangle(.{ .x = x1, .y = y1 }, .{ .x = x2, .y = y2 }, .{ .x = x3, .y = y3 }, s);
}

fn rect(x: i32, y: i32, w: i32, h: i32, s: ff.Style) void {
    ff.drawRect(.{ .x = x, .y = y }, .{ .width = w, .height = h }, s);
}

fn elli(x: i32, y: i32, w: i32, h: i32, s: ff.Style) void {
    ff.drawEllipse(.{ .x = x, .y = y }, .{ .width = w, .height = h }, s);
}

fn circ(x: i32, y: i32, d: i32, s: ff.Style) void {
    ff.drawCircle(.{ .x = x, .y = y }, d, s);
}

fn line(x1: i32, y1: i32, x2: i32, y2: i32, ls: ff.LineStyle) void {
    ff.drawLine(.{ .x = x1, .y = y1 }, .{ .x = x2, .y = y2 }, ls);
}
