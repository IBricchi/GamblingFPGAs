#include <iostream>

using namespace std;

int main(){
    string data;
    cin >> data;
    while(data != "<-->"){
        cin >> data;
    }
    cin >> data;
    while(data != "<-->"){
        cout << data << " ";
        cin >> data;
    }
    cout << endl;
}