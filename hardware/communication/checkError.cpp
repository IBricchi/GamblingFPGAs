#include <iostream>

using namespace std;

int main(){
    string line;
    cin >> line;
    if(line[0] == '{'){
        return 0;
    }
    else
        return 1;
}