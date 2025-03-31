#!/bin/bash

# You need to exchange rsa keys with the server for passwordless uploads
read -p "Push changes to github? (Y\n)" -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]
then
    # YES
    read -p "Enter commit name: " note
    echo
    
    git add -A && git commit -m "$note" ; git push origin master
fi

echo 'Copying css files to remote'
rsync -avh public/css/ aghost:/home/tabula/medical/public/css/
echo 'Done'

echo 'Copying js files to remote'
rsync -avh public/js/ aghost:/home/tabula/medical/public/js/
echo 'Done'

echo 'Copying binary to remote'
rsync -avh medical-go aghost:/home/tabula/medical/medical-go
echo 'Done'

# read -p "Enter password for sudo: " pass

echo 'Restarting remote medical service'
# echo "$pass" | ssh -tt aghost "sudo systemctl restart medical"
ssh -t aghost "sudo systemctl restart medical"
echo 'Done'

