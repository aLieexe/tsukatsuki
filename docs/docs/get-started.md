Tsukatsuki is an open-source command-line tool that streamlines deploying applications on self-hosted infrastructure. It automates the creation of Ansible roles based on user-selected configurations and can also generate GitHub Actions workflows to support CI/CD pipelines. By leveraging Docker as its deployment method, Tsukatsuki helps developers set up consistent, reproducible environments with minimal manual effort.
# Installation
```bash 
git clone https://github.com/aLieexe/tsukatsuki.git
cd tsukatsuki
make install 
```

# Usage
After installing, to start using tsukatsuki, run the following command
```bash
tsukatsuki init
```
Follow the prompt according to your need, tsukatsuki will automatically generate its tsukatsuki.yaml as well as all the file needed to do a deployment
Every configuration is saved in the file `tsukatsuki.yaml`

After generated, run command to deploy to the server
```bash
tsukatsuki deploy
```