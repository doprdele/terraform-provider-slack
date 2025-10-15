module.exports = {
  branches: ['main'],
  plugins: [
    '@semantic-release/commit-analyzer',
    '@semantic-release/release-notes-generator',
    ['@semantic-release/exec', {
      prepareCmd: './build.sh ${nextRelease.version}',
    }],
    ['@semantic-release/github', {
      assets: [
        { path: 'terraform-provider-slack_v*.zip' },
      ]
    }],
    ['@semantic-release/git', {
      assets: ['go.mod', 'go.sum'],
      message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}',
    }],
  ]
};