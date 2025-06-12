#include <iostream>
#include <fstream>
#include <string>
#include <random>
#include <vector>
#include <cmath>
#include <cstring>


/* 

./generate --key_count=100 --read_proportion=0.6 --value_length=4 --distribution=zipf output.txt
./generate --key_count=100 --read_proportion=0.5 --value_length=4 --distribution=uniform output.txt
*/ 
class ZipfianGenerator {
public:
    ZipfianGenerator(int n, double s = 1.0) : N(n), skew(s), dist(0.0, 1.0) {
        normalization_constant = 0.0;
        for (int i = 1; i <= N; ++i) {
            normalization_constant += 1.0 / std::pow(i, skew);
        }
    }

    int next(std::mt19937& gen) {
        double rnd = dist(gen);
        double sum = 0.0;
        for (int i = 1; i <= N; ++i) {
            sum += (1.0 / std::pow(i, skew)) / normalization_constant;
            if (sum >= rnd) return i - 1;
        }
        return N - 1;
    }

private:
    int N;
    double skew;
    double normalization_constant;
    std::uniform_real_distribution<> dist;
};

std::string generate_random_value(int length, std::mt19937& gen) {
    static const char charset[] =
        "0123456789"
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        "abcdefghijklmnopqrstuvwxyz";
    std::uniform_int_distribution<> dist(0, sizeof(charset) - 2);
    std::string result;
    for (int i = 0; i < length; ++i) {
        result += charset[dist(gen)];
    }
    return result;
}

void print_usage(const char* prog_name) {
    std::cerr << "Usage: " << prog_name << " [--key_count=N] [--read_proportion=x] "
              << "[--value_length=L] [--distribution=uniform|zipf] output_file\n";
}

int main(int argc, char* argv[]) {
    int key_count = 10;
    double read_proportion = 0.5;
    int value_length = 4;
    std::string distribution = "uniform";
    std::string output_file;

    // 简单参数解析
    for (int i = 1; i < argc; ++i) {
        if (strncmp(argv[i], "--key_count=", 12) == 0) {
            key_count = std::stoi(argv[i] + 12);
        } else if (strncmp(argv[i], "--read_proportion=", 18) == 0) {
            read_proportion = std::stod(argv[i] + 18);
            if (read_proportion < 0 || read_proportion > 1) {
                std::cerr << "read_proportion must be between 0 and 1\n";
                return 1;
            }
        } else if (strncmp(argv[i], "--value_length=", 15) == 0) {
            value_length = std::stoi(argv[i] + 15);
            if (value_length <= 0) {
                std::cerr << "value_length must be positive\n";
                return 1;
            }
        } else if (strncmp(argv[i], "--distribution=", 15) == 0) {
            distribution = argv[i] + 15;
            if (distribution != "uniform" && distribution != "zipf") {
                std::cerr << "distribution must be 'uniform' or 'zipf'\n";
                return 1;
            }
        } else if (argv[i][0] != '-') {
            output_file = argv[i];
        } else {
            print_usage(argv[0]);
            return 1;
        }
    }

    if (output_file.empty()) {
        print_usage(argv[0]);
        return 1;
    }

    std::ofstream ofs(output_file);
    if (!ofs) {
        std::cerr << "Failed to open output file\n";
        return 1;
    }

    std::random_device rd;
    std::mt19937 gen(rd());

    std::uniform_real_distribution<> op_dist(0.0, 1.0);

    ZipfianGenerator zipf_gen(key_count, 1.2);

    int total_operations = 1000;  // 你可以自己改成参数

    for (int i = 0; i < total_operations; ++i) {
        double op_rnd = op_dist(gen);
        bool is_read = (op_rnd < read_proportion);

        int key_index;
        if (distribution == "uniform") {
            std::uniform_int_distribution<> key_dist(0, key_count - 1);
            key_index = key_dist(gen);
        } else {
            key_index = zipf_gen.next(gen);
        }

        std::string key = "key" + std::to_string(key_index);
        if (is_read) {
            ofs << "./client read " << key << "\n";
        } else {
            std::string value = generate_random_value(value_length, gen);
            ofs << "./client write " << key << " " << value << "\n";
        }
    }

    ofs.close();
    std::cout << "Generated " << total_operations << " operations into " << output_file << "\n";
    return 0;
}

