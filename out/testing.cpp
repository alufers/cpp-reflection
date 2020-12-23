#include <iostream>
#include <iomanip>
#include "sranie.h"
#include <optional>

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
        indent(in);
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
        // indent(in);
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
    if (info->kind == ReflectTypeKind::Vector)
    {
        auto aVec = AnyVectorRef(anything);
        auto size = aVec.size();
        std::cout << "vector " << info->name << ":\n";
        for (size_t i = 0; i < size; i++)
        {
            auto valueAt = aVec.at(i);
            print(valueAt, in + 4);
        }
        return;
    }
    std::cout << "  <unknown>\n";
}

void writeJSONstring(std::ostream &o, const std::string &s)
{
    o << "\"";
    for (auto c = s.cbegin(); c != s.cend(); c++)
    {
        switch (*c)
        {
        case '"':
            o << "\\\"";
            break;
        case '\\':
            o << "\\\\";
            break;
        case '\b':
            o << "\\b";
            break;
        case '\f':
            o << "\\f";
            break;
        case '\n':
            o << "\\n";
            break;
        case '\r':
            o << "\\r";
            break;
        case '\t':
            o << "\\t";
            break;
        default:
            if ('\x00' <= *c && *c <= '\x1f')
            {
                o << "\\u"
                  << std::hex << std::setw(4) << std::setfill('0') << (int)*c;
            }
            else
            {
                o << *c;
            }
        }
    }
    o << "\"";
}

void writeJSON(std::ostream &o, AnyRef any)
{

#define PRINT_LITERALLY(type) \
    if (any.is<type>())       \
    {                         \
        o << *any.as<type>(); \
        return;               \
    }
    PRINT_LITERALLY(bool)
    PRINT_LITERALLY(int)
    PRINT_LITERALLY(double)
    PRINT_LITERALLY(float)
    PRINT_LITERALLY(size_t)
    if (any.is<std::string>())
    {
        writeJSONstring(o, *any.as<std::string>());
        return;
    }
    auto info = any.reflectType();
    if (info->kind == ReflectTypeKind::Class)
    {
        o << "{";
        for (int i = 0; i < info->fields.size(); i++)
        {
            writeJSONstring(o, info->fields[i].name);
            o << ":";
            writeJSON(o, any.getField(i));
            if (i != info->fields.size() - 1)
            {
                o << ",";
            }
        }
        o << "}";
        return;
    }
    if (info->kind == ReflectTypeKind::Enum)
    {
        for (int i = 0; i < info->enumValues.size(); i++)
        {
            int val = *any.as<int>();
            if (info->enumValues[i].value == val)
            {
                writeJSONstring(o, info->enumValues[i].name);
                return;
            }
        }
        return;
    }
    if (info->kind == ReflectTypeKind::Vector)
    {
        auto aVec = AnyVectorRef(any);
        auto size = aVec.size();
        o << "[";
        for (size_t i = 0; i < size; i++)
        {
            auto valueAt = aVec.at(i);
            writeJSON(o, valueAt);
            if (i != size - 1)
            {
                o << ",";
            }
        }
        o << "]";
        return;
    }
}

int main(int argc, char **argv)
{
    std::optional<uint64_t> dziabaDziaba = 5;
    auto optPtr = &dziabaDziaba;
    uint64_t *elemptr = &**optPtr;
    AnyRef refToOptional = AnyRef::of(&dziabaDziaba);
    AnyOptionalRef oref(refToOptional);
    oref.reset();
    std::cout << dziabaDziaba.has_value() << "\n";
    // std::cout << "absolute size of reflection data " << sizeof(reflectTypeInfo) << "\n";
    // Foo foo;
    // foo.alpha = 69;
    // foo.beta = false;
    // foo.gamma = 420;

    // auto a = AnyRef::of(&foo);
    // print(a);

    // std::cout << "\n\n -------\n\n";
    // Bar bar;
    // bar.fooOne = foo;
    // bar.fooTwo.beta = true;
    // bar.fooTwo.gamma = 999;
    // bar.fooTwo.alpha = 777;

    // auto v = AnyRef::of(&bar);
    // print(v);

    // std::cout << "\n\n -------\n\n";

    // auto theTypeOfBar = AnyRef::of(v.reflectType());
    // print(theTypeOfBar);

    Frame baz;
    // baz.foos.push_back(foo);
    // baz.foos.push_back(foo);
    auto bazRef = AnyRef::of(&baz);
    auto theTypeOfBaz = AnyRef::of(bazRef.reflectType());
    writeJSON(std::cout, theTypeOfBaz);

    // std::cout << "\n\n -------\n\n";

    // std::cout << "THE ABSOLUTE TYPE OF BAZ:";

    // std::cout << "\n\n -------\n\n";
    // auto theTypeOfBaz = AnyRef::of(bazRef.reflectType());
    // print(theTypeOfBaz);

    // auto created = AnyRef::construct<Foo2>();
    // print(created);
    return 0;
}
