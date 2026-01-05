const std = @import("std");

pub fn build(b: *std.Build) !void {
    const id = "ff-author-id.ff-app-id";

    const exe = b.addExecutable(.{
        .name = id,
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/main.zig"),
            .target = b.resolveTargetQuery(.{
                .cpu_arch = .wasm32,
                .os_tag = .freestanding,
            }),
            .optimize = .ReleaseSmall,
        }),
    });

    exe.root_module.addImport("ff", b.dependency("ff", .{}).module("ff"));

    exe.entry = .disabled;
    exe.rdynamic = true;

    b.installArtifact(exe);

    const build_cmd = b.addSystemCommand(&[_][]const u8{
        "firefly_cli",
        "build",
    });

    const run_cmd = b.addSystemCommand(&[_][]const u8{
        "firefly-emulator",
        "--id",
        id,
    });
    run_cmd.step.dependOn(&build_cmd.step);

    const run_step = b.step("run", "Run app in the Firefly Zero emulator");

    run_step.dependOn(&run_cmd.step);

    const spy_cmd = b.addSystemCommand(&[_][]const u8{
        "spy",
        "--exc",
        "zig-cache",
        "--inc",
        "**/*.zig",
        "-q",
        "./spy.sh",
    });

    const spy_step = b.step("spy", "Run spy watching for file changes");
    spy_step.dependOn(&spy_cmd.step);
}
