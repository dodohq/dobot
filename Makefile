setup:
	rm -f .git/hooks/pre-commit.sample
	curl https://gist.githubusercontent.com/proxypoke/2781569/raw/091ca83dc9c1e460009a3e4ed8988a5803a423dd/pre-commit.sh > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

dev:
	GO_ENV=development fresh