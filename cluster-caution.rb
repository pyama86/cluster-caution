# This file was generated by GoReleaser. DO NOT EDIT.
class ClusterCaution < Formula
  desc "Prevents execution errors in kubectl"
  homepage "https://github.com/pyama86/cluster-caution"
  url "https://github.com/pyama86/cluster-caution/releases/download/0.1.0/cluster-caution_0.1.0_darwin_amd64.tar.gz"
  version "0.1.0"
  sha256 "a873b59fce14811b2d5a8560d0df862c15d925dd7332a40a5ba9fc556727930f"

  def install
    bin.install 'kubectl-cluster-caution'
  end

  test do
    system "#{bin}/kubectl-cluster-caution"
  end
end
