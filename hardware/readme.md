## Quartus Setup

The recomended setup for this project follows the following steps:

1. From https://fpgasoftware.intel.com/20.1.1/?edition=lite&platform=windows download:
    - Quartus Prime (includes Nios II EDS)
    - ModelSim-Intel FPGA Edition (includes Starter Edition)
    - MAX 10 FPGA device support

    Shuld work with versions of quartus after 19, and with modifications could potentially work with older versions, but not tested
2. For versions of quartus after 19 eclipse must be installed manually. Follow [instructions](https://www.intel.com/content/altera-www/global/en_us/index/support/support-resources/knowledge-base/tools/2019/why-does-the-nios--ii-not-installed-after-full-installation-of-t.html) for intell install guide. Very well explained.
3. This requires Ubuntu 18.04 installed on wsl. [Instructions on setting up wsl](https://docs.microsoft.com/en-us/windows/wsl/install-win10). You probably aren't a windows insider, so scroll down a little to find the normall method.
4. Download Ubuntu 18.04 from the Microsoft Store
5. Open a power shell terminal and execute the following command
    - ``` wsl -l ```
    - Copy the version of that either sais Ubuntu or Ubuntu-18.04
    - ``` wsl --set-version 1 ```
    - nios2 needs wsl version 1
    - ``` wsl -s <Name of Distro> ```
    - nios2 uses the default distro
    - ``` wsl ```
    - enter wsl terminal
    - ``` sudo apt update ```
    - ``` sudo apt-get install wsl ```
    - ``` sudo apt-get install g++ ```
    - ``` sudo apt install dos2unix ```
    - ``` sudo apt install make ```
    - nios2 uses these packages
6. The final step is to edit your path:
    - Type in 'Edit environment variables' in the start menu, and open the control panel program
    - In the list select 'Path' and click the edit button
    - On the new window click new, and add the following line:
        ``` C:\intelFPGA_lite\20.1\quartus\bin64 ```

Everythign should now be set up and working.

## Build project

Only the minimum required files are included in the repo so to set everything up a few files need to be compiled on your system.

1. Open Quartus, select open project, and select the bdf file in /hardware/FPGA
2. This should set everything up correctly, no need to add extra files or anything
2. In the files tab on quartus you should see a qsys file, double click on that, and compile to HDL
3. After this you can now compile you're entire project, this takes a while but only needs to be done once so at least there's that
4. From here you can program you device
4. Open eclipse tools for nios2
5. For the workspace I recomend selecting /hardware/workspace it has .gitignore setup to ignore almost everything
6. Copy the sopcinfo file from /hardware/FPGA/ into the workspace folder
7. Create new 'NIOS II Application and BSP from Template'
    - Select SOPC file to be the sopcinfo file you copied
    - Set name to whatever you want (setting it to fpga is recomeded)
    - From project template select hello Wolrd
    - Delete the hello_world.c file
    - From file explorer on windows, drag both main.c and the src folder from /hardware/workspace into your created project on the eclipse side bar. IT IS IMPORTANT YOU DRAG IT INTO THE ECLIPSE SIDEBAR so that the files get properly linked up. If you see a popup select copy files the other options will not work. This should link up the make files appropriately and update the makefiles so eclipse can handle building.

## Developing

Most files are ignored by git as to keep the repo size small. If big changes are made to the quartus file that creates files that are currently ignored, try and find out what the minimum number of files necessary are to generate the whole project, and ammend the build project instructions if new steps are required.

As for the c program, changes in your workspace will not be caputred by git, so remember to copy the contents of the main.c and src from your project back into the file /hardware/workspace/. If you've setup your workspace as recomended in the setup steps you can run ./projtosrc.sh to automate this.
### Warning

Just in case you accidentally use the projtosrc.sh command at the wrong time. The file makes a backup of the main.c file and src folder in /hardware/workspace/. This only caputre the perevious version of these files so be careful not to accidentally run the helper tool twice.

## Run Project

1. On quartus compile, and programme the fpga
2. Before running the eclipse project open an instance of powershell and run the command ```nios2-terminal.exe```
2. MAKE SURE YOU DID THE PREVIOUS STEP AS IF YOU DON'T AND RUN THE PROJECT, ECLIPSE WILL TAKE OVER YTHE NIOS2-TERMINAL AND YOU WILL HAVE TO UNPLUG THE FPGA AND RE-PROGRAMME IT!!!!
2. On eclipse (or using the cli for eclipse) compile and run the project.
2. Once the program is running you should see the message "Running.." written on the powershell you opened
3. You can now press CTR-C to terminate the nios2-terminal on powershell, DO NOT PRESS CTR-Z as this sometimes terminates the C programme as well
3. Open wsl (or got through powershell using the the ```wsl``` command), go to /hardware/communication folder
4. Run ./setup.sh to compiler any requirements to run the proper script running
5. Run ```./requestLoop.sh <username> <password>```
6. If you get any errors in about non existing files and the shell scripts do not run, try using ```dos2unix``` on the shell scripts and re-run.
7. requestLoop.sh should print some messages about download speed, one json file, another json file, and then repeat that over and over again. If anywhere in the output you see a string that looks like it could say 'terminate' but is weridly chopped off, as long as it only happens occasionally that is expected. If it's happening all the time then unplug your FPGA and re-programme it. If that doesn't solve the problem, or if you're getting any other wiered messages please let me know (Ignacio) and we can try and figure out why.
