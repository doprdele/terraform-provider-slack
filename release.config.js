module.exports = {
  branches: ['main'],
  tagFormat: 'v${version}',
  plugins: [
    '@semantic-release/commit-analyzer',
    '@semantic-release/release-notes-generator',
    ['@semantic-release/exec', {
      prepareCmd: './build.sh ${nextRelease.version}',
    }],
    ['@semantic-release/github', {
      assets: [
        { path: 'terraform-provider-slack_*.zip' },
        { path: 'terraform-provider-slack_*_SHA256SUMS' },
        { path: 'terraform-provider-slack_*_SHA256SUMS.sig' },
      ]
    }],
    ['@semantic-release/git', {
      assets: ['go.mod', 'go.sum'],
      message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}',
    }],
  ]
};
