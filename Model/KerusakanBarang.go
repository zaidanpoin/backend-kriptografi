package Model

import "kriptografi-zaidaan/Database"

type KerusakanBarang struct {
	KodeBarang string `json:"kode_barang" gorm:"primaryKey"`
	Deskripsi  string `json:"deskripsi"`
	Tanggal    string `json:"tanggal"`
	Barang     Barang `json:"barang" gorm:"foreignKey:KodeBarang"`
}

func (k *KerusakanBarang) Save() (*KerusakanBarang, error) {
	if err := Database.Database.Create(k).Error; err != nil {
		return nil, err
	}
	return k, nil
}
func (k *KerusakanBarang) GetKerusakanByKode(kodeBarang string) (*KerusakanBarang, error) {
	var kerusakan KerusakanBarang
	if err := Database.Database.Preload("Barang").First(&kerusakan, "kode_barang = ?", kodeBarang).Error; err != nil {
		return nil, err
	}
	return &kerusakan, nil
}

func (k *KerusakanBarang) GetAllKerusakan() ([]KerusakanBarang, error) {
	var kerusakans []KerusakanBarang
	if err := Database.Database.Preload("Barang").Find(&kerusakans).Error; err != nil {
		return kerusakans, err
	}
	return kerusakans, nil
}
func (k *KerusakanBarang) DeleteKerusakan(kodeBarang string) error {
	if err := Database.Database.Delete(&KerusakanBarang{}, "kode_barang = ?", kodeBarang).Error; err != nil {
		return err
	}
	return nil
}
func (k *KerusakanBarang) UpdateKerusakan(kodeBarang string, updatedKerusakan KerusakanBarang) (*KerusakanBarang, error) {
	var kerusakan KerusakanBarang
	if err := Database.Database.First(&kerusakan, "kode_barang = ?", kodeBarang).Error; err != nil {
		return nil, err
	}

	kerusakan.Deskripsi = updatedKerusakan.Deskripsi
	kerusakan.Tanggal = updatedKerusakan.Tanggal

	if err := Database.Database.Save(&kerusakan).Error; err != nil {
		return nil, err
	}
	return &kerusakan, nil
}
