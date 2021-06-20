alias dev-proxy-staging='docker run -it --rm -p 80:80 -p 8182:8182 -p 4444:4444 framework-s-staging.tu-server-slv01.corp.waters.com/dev-proxy:latest --cluster scalars-dev --user ukafig --os win' 
alias dev-proxy='docker run -it --rm -p 80:80 -p 8182:8182 -p 4444:4444 032356282346.dkr.ecr.eu-west-2.amazonaws.com/dev-proxy:latest --cluster scalars-dev --user ukafig --os win' 
alias aws-docker-login='REGION=eu-west-2 && aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin 032356282346.dkr.ecr.${REGION}.amazonaws.com'

alias config='/usr/bin/git --git-dir=/home/adam/.myconfig/ --work-tree=/home/adam'
