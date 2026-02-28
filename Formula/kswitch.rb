class Kswitch < Formula
  desc "Interactive Kubernetes context switcher with arrow-key navigation"
  homepage "https://github.com/YonierGomez/kswitch"
  url "https://github.com/YonierGomez/kswitch/archive/refs/tags/v1.2.0.tar.gz"
  sha256 "1001a75f2b5b266e5157c9710863973141d31d59ee19868f65af162d130714f3"
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
