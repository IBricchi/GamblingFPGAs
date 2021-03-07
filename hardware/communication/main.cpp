#include <iostream>

using namespace std;

int main(){
    string line;
    cin >> line;
    while(!cin.fail()){
        if(line == "<data>") break;
        cin >> line;
    }
    cin >> line;
    cout << line << endl;
    return 0;
}