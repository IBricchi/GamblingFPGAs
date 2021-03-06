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

1. Open quartus, select open project, and select the bdf file in /hardware/Golden_Top
2. In files you should see a qsys file, double click on that, and compile to HDL
3. From here you can program you device
4. Open eclipse tools for nios2
5. You'll have to create a workspace
    - I recomened making a folder called workspace in /hardware, it will be fully ignored by git
6. Copy the socpinfo file from /hardware/Golden_Top/nios_accelerometer/ into your workspace
7. Create new 'NIOS II Application and BSP from Template'
    - Select SOPC file to be the sopcinfo file you copied
    - Select name to fpga or whatever you prefer (will be ignored in git anyway)
    - From project template select hello Wolrd
    - Replace the hello_world.c file in your workspace with the one from /hardware/Golden_Top/hello_world.c 

Now you should be able to run the project

## Developing

Most files are ignored by git as to keep the repo size small. If big changes are made to the quartus file that creates files that are currently ignored, try and find out what the minimum number of files necessary are to generate the whole project, and ammend the build project instructions if new steps are required.

As for the c program, changes in your workspace will not be caputred by git, so remember to always overwire the hello_world.c file in /hardware/

## Run Project

1. On quartus compile, and programme the fpga
2. On eclipse (or using the cli for eclipse) compile and run the project.
3. Open wsl (or powershell and then the wsl command), go to /hardware/communication folder
4. Run ./run <char>, and you should see a response printed onto your screen from the fpga 
5. Be careful because if main.cpp doesn't recieve expected sygnal it enters infinite loop, this is only a quick test I came up with, will fix this later