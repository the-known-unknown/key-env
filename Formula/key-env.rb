class KeyEnv < Formula
  desc "Load env vars and hydrate secrets from local vault CLIs"
  homepage "https://github.com/asimmittal/key-env"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/asimmittal/key-env/releases/download/v#{version}/key-env_#{version}_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_ARM64_SHA256"
    else
      url "https://github.com/asimmittal/key-env/releases/download/v#{version}/key-env_#{version}_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_AMD64_SHA256"
    end
  end

  depends_on "keepassxc"

  def install
    bin.install "key-env"
  end

  test do
    assert_match "usage", shell_output("#{bin}/key-env 2>&1", 1)
  end
end
