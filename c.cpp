// C++ example

#include<iostream>

unsigned short arr[] = {1000, 2000, 3000};

struct ST {
	unsigned short rptr :6;
	unsigned short wptr :6;
	unsigned short seq  :2;
};

int func1__() {
	std::cout << "func()" << std::endl;

	return 999;
}

void func1() {
	std::cout << "Hello World!" << std::endl;

	static ST st;


	std::cout <<  arr << std::endl;
	std::cout << &arr << std::endl;
	std::cout << *arr << std::endl;
	std::cout << *arr << std::endl;


	std::cout <<  func1__ << std::endl;
	std::cout << &func1__ << std::endl;
	std::cout << func1__() << std::endl;
	std::cout << (*func1__)() << std::endl;
	std::cout << (********************************************func1__)() << std::endl;
}

void func2() {
	unsigned short s = 0xFFFF;
	unsigned int   t = 0xFFFF;

	unsigned long long  ch = 16;

	std::cout << std::hex << sizeof (s << ch) << std::endl;
	std::cout << std::hex << sizeof (t << ch)<< std::endl;

	s <<= 16;
	t <<= 16;
	std::cout << std::hex << s << std::endl;
	std::cout << std::hex << t << std::endl;
}

void func3() {
	std::cout << (1 /2 + 1)<< std::endl;
	std::cout << (1 /2 + 1.1)<< std::endl;
	std::cout << (1.1 + 1.0 / 2)<< std::endl;
}
int main() {
	//func1();
	func2();
	//func3();
}