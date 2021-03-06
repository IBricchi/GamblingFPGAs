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
    - ``` sudo apt-get install wsl ```
    - nios2 uses this package
6. The final step is to edit your path:
    - Type in 'Edit environment variables' in the start menu, and open the control panel program
    - In the list select 'Path' and click the edit button
    - On the new window click new, and add the following line:
        ``` C:\intelFPGA_lite\20.1\quartus\bin64 ```

Everythign should now be set up and working.

## Build project

Only the minimum required files are included in the repo so to set everything up a few files need to be compiled on your system.

1. Open Quartus, select open project, and select the bdf file in /hardware/FPGA
2. In files you should see a qsys file, double click on that, and compile to HDL
3. From here you can program you device
4. Open eclipse tools for nios2
5. For the workspace I recomend selecting /hardware/workspace it has .gitignore setup to ignore almost everything
6. Copy the sopcinfo file from /hardware/FPGA/ into the workspace folder
7. Create new 'NIOS II Application and BSP from Template'
    - Select SOPC file to be the sopcinfo file you copied
    - Set name to whatever you want (setting it to fpga is recomeded)
    - From project template select hello Wolrd
    - Copy the contents of hardware/workspace/hello_world.c into the hello_world.c file in your eclipse project. (If you've set it up as recomended, with the project name fpga, you can run ./workspace/srctoproj.sh to automate this)

Now you should be able to run the project

## Developing

Most files are ignored by git as to keep the repo size small. If big changes are made to the quartus file that creates files that are currently ignored, try and find out what the minimum number of files necessary are to generate the whole project, and ammend the build project instructions if new steps are required.

As for the c program, changes in your workspace will not be caputred by git, so remember to copy the contents of the hello_world.c file from your project back into the file /hardware/workspace/hello_world.c. If you've setup your workspace as recomended in the setup steps you can run ./workspace/projtosrc.sh to automate this.

### Warning

Just in case you accidentally use the wrong one of the projtosrc or srctoproj, both commands make a backup of the src and proj files which can be found in hardware/workspace/bak/ these only hold the state of both files before the last command, so if you run any of the commands twice you could loose changes.

## Run Project

1. On quartus compile, and programme the fpga
2. On eclipse (or using the cli for eclipse) compile and run the project.
3. Open wsl (or powershell and then the wsl command), go to /hardware/communication folder
4. Run ./run.sh <char>, and you should see a response printed onto your screen from the fpga 
5. Be careful because if main.cpp doesn't recieve expected sygnal it enters infinite loop, this is only a quick test I came up with, will fix this later