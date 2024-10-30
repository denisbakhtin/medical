#!/bin/bash

# You need to exchange rsa keys with the server for passwordless uploads
read -p "Push changes to github? (Y\n)" -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]
then
    # YES
    read -p "Enter commit name: " -r
    echo
    
    git add -A && git commit -m "$REPLY" ; git push origin master
fi

# echo 'Copying css files to remote'
# rsync -avh public/css/ aghost:/home/tabula/medical/public/css/
# echo 'Done'
#
# echo 'Copying js files to remote'
# rsync -avh public/js/ aghost:/home/tabula/medical/public/js/
# echo 'Done'
#
# echo 'Copying binary to remote'
# rsync -avh miobalans-go aghost:/home/tabula/medical/miobalans-go
# echo 'Done'
#
# echo 'Restarting remote medical service'
# ssh -t aghost "sudo systemctl restart medical"
# echo 'Done'
