#include <iostream>
#include <iomanip>
#include "sranie.h"

void indent(int n)
{
    for (int i = 0; i < n; i++)
    {
        std::cout << " ";
    }
}

void print(AnyRef &anything, int in = 0)
{

    if (anything.is<bool>())
    {
        bool val = *anything.as<bool>();
        std::cout << "boolean: " << val << "\n";
        return;
    }
    else if (anything.is<int>())
    {
        int val = *anything.as<int>();
        std::cout << "int: " << val << "\n";
        return;
    }
    else if (anything.is<size_t>())
    {
        size_t val = *anything.as<size_t>();
        std::cout << "size_t: " << val << "\n";
        return;
    }
     else if (anything.is<std::string>())
    {
        std::string val = *anything.as<std::string>();
        std::cout << "std::string: '" << val << "'\n";
        return;
    }


    auto info = anything.reflectType();
    if (info->kind == ReflectTypeKind::Class)
    {

        int val = *anything.as<int>();

        
        std::cout << "class " << info->name << ":\n";
        for (int i = 0; i < info->fields.size(); i++)
        {
            indent(in + 4);
            std::cout << std::left << std::setw(7) << info->fields[i].name << " ";
            auto fieldAny = anything.getField(i);
            print(fieldAny, in + 8);
        }
        return;
    }
    if (info->kind == ReflectTypeKind::Enum)
    {

        int val = *anything.as<int>();

       
        std::cout << "enum " << info->name << ": ";
        for (int i = 0; i < info->enumValues.size(); i++)
        {
            if (info->enumValues[i].value == val)
            {
                std::cout << info->enumValues[i].name;
            }
        }
        std::cout << "\n";
        return;
    }
    std::cout << "  <unknown>\n";
}

int main(int argc, char **argv)
{
    std::cout << "size of reflection data " << sizeof(std::string) << "\n";
    Foo foo;
    foo.alpha = 69;
    foo.beta = false;
    foo.gamma = 420;

    auto a = AnyRef::of(&foo);
    print(a);

    std::cout << "\n\n -------\n\n";
    Bar bar;
    bar.fooOne = foo;
    bar.fooTwo.beta = true;
    bar.fooTwo.gamma = 999;
    bar.fooTwo.alpha = 777;

    auto v = AnyRef::of(&bar);
    print(v);

    std::cout << "\n\n -------\n\n";

    auto theTypeOfBar = AnyRef::of(v.reflectType());
    print(theTypeOfBar);

    // auto created = AnyRef::construct<Foo2>();
    // print(created);
    return 0;
}
