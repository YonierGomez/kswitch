class Kswitch < Formula
  desc "Interactive Kubernetes context switcher with arrow-key navigation"
  homepage "https://github.com/YonierGomez/kswitch"
  url "https://github.com/YonierGomez/kswitch/archive/refs/tags/v1.1.1.tar.gz"
  sha256 "9991d863ca9e5c19af070969a1ae045fb916987aa2001452c8316646e916b0a2"
  license "MIT"

  depends_on "go" => :build
  depends_on "kubernetes-cli"

  def install
    system "go", "build", "-ldflags", "-s -w", "-o", bin/"ksw", "."
  end

  test do
    assert_match "ksw v#{version}", shell_output("#{bin}/ksw -v")
  end
end
