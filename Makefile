# START: begin
GITHUB_BACKUP_PATH=${HOME}/go/src/github.com/srinathLN7/gitbackup/gossip-gloomers


.PHONY: gitbackup
gitbackup:
	echo "backing up .git dir and pushing code to `main` branch"
	sudo cp -rf ./.git  ${GITHUB_BACKUP_PATH}
	git push
	
# END: begin
