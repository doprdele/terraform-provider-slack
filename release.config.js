module.exports = {
  branches: ['main'],
  plugins: [
    '@semantic-release/commit-analyzer',
    '@semantic-release/release-notes-generator',
    ['@semantic-release/exec', {
      prepareCmd: './build.sh ${nextRelease.version} && gpg --clearsign --armor --batch --yes --output SHA256SUMS.sig SHA256SUMS',
    }],
    ['@semantic-release/github', {
      assets: [
        { path: 'terraform-provider-slack_v*.zip' },
        { path: 'SHA256SUMS' },
        { path: 'SHA256SUMS.sig' },
      ]
    }],
    ['@semantic-release/git', {
      assets: ['go.mod', 'go.sum'],
      message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}',
    }],
  ]
};