# Homebrew formula for MinifyJS.
#
# This formula is a template. The `url` and `sha256` fields are
# filled in by the release pipeline (build/generate_checksums.py
# emits a version suitable for pasting here). The `version` field
# must match the release version.
#
# To submit the formula upstream, open a pull request against
# https://github.com/Homebrew/homebrew-core with this file, with
# the placeholder values replaced by the actual release values.
#
# Users who do not want to wait for upstream submission can install
# from a tap:
#
#     brew tap minifyjs/minifyjs
#     brew install minifyjs
#
# The tap repository mirrors this formula with real values.

class Minifyjs < Formula
  desc "A native JavaScript minifier and optimizer (no Node.js required)"
  homepage "https://github.com/minifyjs/minifyjs"
  version "0.1.0"
  license "MIT"

  # The formula supports both Intel and Apple Silicon. Homebrew
  # picks the right one based on the user's architecture.
  on_macos do
    on_intel do
      url "https://github.com/minifyjs/minifyjs/releases/download/v#{version}/minifyjs-darwin-amd64"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end

    on_arm do
      url "https://github.com/minifyjs/minifyjs/releases/download/v#{version}/minifyjs-darwin-arm64"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end
  end

  # Linux users can install with:
  #
  #     brew install minifyjs
  #
  # when running Homebrew on Linux. The formula picks the right
  # binary for the architecture.
  on_linux do
    on_intel do
      url "https://github.com/minifyjs/minifyjs/releases/download/v#{version}/minifyjs-linux-amd64"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end

    on_arm do
      url "https://github.com/minifyjs/minifyjs/releases/download/v#{version}/minifyjs-linux-arm64"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end
  end

  def install
    # The downloaded file is the raw binary. Homebrew saves it as
    # "minifyjs" (the version is stripped from the filename), but
    # to be safe we look for either name.
    binary = Dir["minifyjs*"].first
    raise "binary not found in archive" unless binary

    bin.install binary => "minifyjs"
  end

  def caveats
    <<~EOS
      MinifyJS is a JavaScript minifier that runs without Node.js.

      Try it:
        echo "const x = 1 + 2 + 3;" | minifyjs --compress
        # const x=6;

      Full documentation:
        https://github.com/minifyjs/minifyjs/tree/main/docs
    EOS
  end

  test do
    # Homebrew runs this after install to verify the formula works.
    assert_match "minifyjs", shell_output("#{bin}/minifyjs --version")
    assert_match "esbuild", shell_output("#{bin}/minifyjs --version")

    assert_equal(
      "const x=1;",
      pipe_output("#{bin}/minifyjs", "const x = 1;\n"),
    )

    assert_equal(
      "const x=6;",
      pipe_output("#{bin}/minifyjs --compress", "const x = 1 + 2 + 3;\n"),
    )
  end
end