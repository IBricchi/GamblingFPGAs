#include <iostream>
#include <string>

using namespace std;

int main(){
    string line;
    getline(cin, line);
    while(!cin.fail()){
        if(line == "<data>") break;
        getline(cin, line);
    }
    cin >> line;
    cout << line << endl;
    return 0;
}